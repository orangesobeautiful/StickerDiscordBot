from Database.SQLAlchemyDrawImageOperation import SQLAlchemyDrawImageOperation
from GoogleDriverAPI import GDBaseOp
import os
import operator
import pathlib
import time


def cmp(a, b):
    return (a > b) - (a < b)


class DrawImage(GDBaseOp.GDBaseOp):

    def __init__(self):
        self.folder_image_num = 200
        super(DrawImage, self).__init__()
        self.db_url = os.environ['DATABASE_URL']
        self.db_op = SQLAlchemyDrawImageOperation(self.db_url)

    def add_image_source(self):
        folder_id, path = self.select_folder()
        if self.db_op.add_image_source(folder_id, path):
            print('{0} 新增成功'.format(path))
        else:
            print('{0} 新增失敗'.format(path))

    def all_image_source(self, show: bool):
        count = 0
        image_source_list = self.db_op.all_image_source()
        if show:
            for source in image_source_list:
                print('{index:03d}. {folder_id:s}\t({path:s})'.format(index=count, path=source[1], folder_id=source[0]))
                count = count + 1

        return image_source_list

    def remove_source(self):
        selected = False
        while not selected:
            print('All Image Source:')
            image_source_list = self.all_image_source(show=True)
            input_str = input('Source編號:')
            try:
                source_num = int(input_str)
                self.db_op.delete_source(image_source_list[source_num][0])
            except ValueError:
                if input_str.lower() == 'exit':
                    break
                else:
                    print('請輸入數字')

    def update_images(self):
        image_source_list = self.all_image_source(show=False)

        for image_source in image_source_list:
            source_folder_id = image_source[0]
            source_folder_path = image_source[1]

            all_folder_list: list = self.get_folders(parent_id=source_folder_id)
            updated_folder_list: list = self.db_op.all_updated_folders(parent_folder=source_folder_id)
            all_folder_list.sort(key=lambda x: x['id'])
            updated_folder_list.sort(key=lambda x: x[0])
            len_updated_folder_list = len(updated_folder_list)
            len_all_folder_list = len(all_folder_list)
            none_updated_folder = list()

            if len_updated_folder_list == 0:
                none_updated_folder = all_folder_list
            elif len_all_folder_list == 0:
                print('{source_path:s} 裡面沒資料夾'.format(source_path=source_folder_path))
            else:
                i = j = 0
                while True:
                    if i >= len_all_folder_list:
                        i = len_all_folder_list - 1
                    if j >= len_updated_folder_list:
                        j = len_updated_folder_list - 1

                    id_cmp = cmp(all_folder_list[i]['id'], updated_folder_list[j][0])
                    if id_cmp > 0:
                        none_updated_folder.append(all_folder_list[i])
                        i += 1
                        j += 1
                    elif id_cmp < 0:
                        none_updated_folder.append(all_folder_list[i])
                        i += 1
                    else:
                        j += 1
                        i += 1

                    if i >= len_all_folder_list and j >= len_updated_folder_list:
                        break
            # print(len(none_updated_folder))
            for folder in none_updated_folder:
                folder_id = folder['id']
                folder_path = os.path.join(source_folder_path, folder['name'])
                all_items_list = self.get_all_items(folder_id)
                if len(all_items_list) >= self.folder_image_num:
                    for item in all_items_list:
                        self.db_op.add_images(item['id'], pathlib.Path(os.path.join(folder_path, item['name'])).as_posix())
                    self.db_op.add_updated_folder(folder_id, pathlib.Path(folder_path).as_posix(), source_folder_id)
                    print('{folder_path} 新增完畢'.format(folder_path=pathlib.Path(folder_path).as_posix()))

    def test(self):
        image_list = self.db_op.all_images()
        count = 0
        pre_folder = image_list[0][1][:image_list[0][1].rfind('/')]
        for image in image_list:
            last_index = image[1].rfind('/')
            cur_folder = image[1][:last_index]
            if cur_folder != pre_folder:
                print('{folder_path}:{num}'.format(folder_path=pre_folder, num=count))
                pre_folder = cur_folder
                count = 0
            count += 1

    def get_rand_image(self):
        image_id = self.db_op.get_rand_image()
        # print(image_id)
        if image_id:
            image_url = 'https://drive.google.com/uc?export=view&id=' + image_id[0]
            return image_url
        else:
            return None

    def get_rand_image_id(self):
        image_id = self.db_op.get_rand_image()
        # print(image_id)
        if image_id:
            return image_id[0]
        else:
            return None

    def clear_data(self):
        self.db_op.delete_all_image()
        self.db_op.delete_updated_folders()


if __name__ == '__main__':
    gd = DrawImage()
    print(gd.get_rand_image_id())