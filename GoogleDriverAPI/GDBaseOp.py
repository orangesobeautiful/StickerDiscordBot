import pickle
import os.path
from googleapiclient.discovery import build
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request
import os
import pathlib

# If modifying these scopes, delete the file token.pickle.
SCOPES = ['https://www.googleapis.com/auth/drive.metadata.readonly']


class GDBaseOp:

    def __init__(self):
        self.service = None
        self.get_cred()

    def get_cred(self):
        """Shows basic usage of the Drive v3 API.
        Prints the names and ids of the first 10 files the user has access to.
        """
        creds = None
        # The file token.pickle stores the user's access and refresh tokens, and is
        # created automatically when the authorization flow completes for the first
        # time.
        if os.path.exists('token.pickle'):
            with open('token.pickle', 'rb') as token:
                creds = pickle.load(token)
        # If there are no (valid) credentials available, let the user log in.
        if not creds or not creds.valid:
            if creds and creds.expired and creds.refresh_token:
                creds.refresh(Request())
            else:
                flow = InstalledAppFlow.from_client_secrets_file(
                    'credentials.json', SCOPES)
                creds = flow.run_console()
                #print(authorize_url)
                #creds = flow.run_local_server(port=0)
            # Save the credentials for the next run
            with open('token.pickle', 'wb') as token:
                pickle.dump(creds, token)

        self.service = build('drive', 'v3', credentials=creds)
        return self.service

    def select_folder(self):
        # Call the Drive v3 API
        selected = False
        folder_name = 'root'
        folder_id = 'root'
        cur_path = 'root'
        up_id = 'root'
        pre_folder = 'root'
        while not selected:
            query = "'{fid}' in parents and mimeType = 'application/vnd.google-apps.folder'".format(fid=folder_id)
            results = self.service.files().list(q=query, pageSize=100, spaces='drive',
                                                fields="nextPageToken, files(name, id)", orderBy='name').execute()
            items = results.get('files', [])

            print('\nFolders({path}):'.format(path=cur_path))
            if items:
                for i in range(len(items)):
                    item = items[i]
                    print(u'{i:03d}.{name}'.format(i=i, name=item['name']))
            else:
                print('No items')

            input_str = input('資料夾編號(輸入Select選定):')
            try:
                folder_num = int(input_str)
                up_id = folder_id
                folder_name = items[folder_num]['name']
                folder_id = items[folder_num]['id']
                cur_path = os.path.join(cur_path, folder_name)
                cur_path = pathlib.Path(cur_path).as_posix()

            except ValueError:
                if input_str.strip().lower() == 'root':
                    up_id = 'root'
                    folder_id = 'root'
                    cur_path = 'root'
                elif 'select' in input_str.strip().lower():
                    selected = True
                elif 'up' in input_str.strip().lower():
                    folder_id = up_id
                    if not cur_path == 'root':
                        cur_path = os.path.dirname(cur_path)
                        cur_path = pathlib.Path(cur_path).as_posix()
                else:
                    print('錯誤的輸入')

        return folder_id, cur_path

    def get_folders(self, parent_id):
        all_file = list()
        next_page_token = None
        while True:
            query = "'{fid}' in parents and mimeType = 'application/vnd.google-apps.folder'".format(fid=parent_id)
            results = self.service.files().list(q=query, pageSize=200, spaces='drive',
                                                fields="nextPageToken, files(name, id)", orderBy='name',
                                                pageToken=next_page_token).execute()
            items = results.get('files')
            all_file.extend(items)
            next_page_token = results.get('nextPageToken')
            if next_page_token is None:
                break
        return all_file

    def get_all_items(self, parent_id):
        all_item = list()
        next_page_token = None
        while True:
            query = "'{fid}' in parents".format(fid=parent_id)
            results = self.service.files().list(q=query, pageSize=200, spaces='drive',
                                                fields="nextPageToken, files(name, id)", orderBy='name',
                                                pageToken=next_page_token).execute()
            items = results.get('files')
            all_item.extend(items)
            next_page_token = results.get('nextPageToken')
            if next_page_token is None:
                break
        return all_item


if __name__ == '__main__':
    gd = GDBaseOp()
    print(gd.select_folder())
    #for e in gd.get_folders('root'):
    #    print(e)
    'https://drive.google.com/uc?export=view&id=1cU9ZSVpKig5AOFoEiavsgdqZf-LoVc1w'
