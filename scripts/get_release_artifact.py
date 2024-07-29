import requests
import json
import os
import re
import sys
import shutil


def setup_artifacts_folder(folder_path):
    if os.path.exists(folder_path):
        shutil.rmtree(folder_path)
    os.makedirs(folder_path)

def download_latest_release(repository,base_version):
    pat = os.getenv('PAT')

    if not pat:
        raise ValueError("Please set the PAT environment variable.")

    url = f"https://api.github.com/repos/icon-project/{repository}/releases"
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": f"Bearer {pat}",
        "X-GitHub-Api-Version": "2022-11-28"
    }

    response = requests.get(url, headers=headers)

    pattern = re.compile(rf'^(v)?{base_version}(?:-[a-zA-Z0-9._-]+)?(-hotfix)?$')

    artifacts_folder = "artifacts"
    setup_artifacts_folder(artifacts_folder)

    try:
        data = response.json()
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        data = []

    if isinstance(data, list):
        latest_release = None
        latest_tag = None

        for release in data:
            if isinstance(release, dict):
                if not release.get('draft', False) and not release.get('prerelease', False):
                    tag_name = release.get('tag_name')
                    if tag_name and pattern.match(tag_name):
                        print(f"Found valid release: {tag_name}")

                        assets = release.get('assets', [])
                        for asset in assets:
                            if isinstance(asset, dict):
                                download_url = asset.get('browser_download_url')
                                if download_url:
                                    response = requests.get(download_url)
                                    file_name = download_url.split('/')[-1]
                                    file_path = os.path.join(artifacts_folder, file_name)
                                    with open(file_path, 'wb') as file:
                                        file.write(response.content)
                                    print(f"Downloaded {file_name} from {download_url}")

    else:
        print("Unexpected data format:", type(data))

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python get_release_artifact.py <repository> <base_version>")
        sys.exit(1)

    repository_name = sys.argv[1]
    base_version = sys.argv[2]
    download_latest_release(repository_name, base_version)
