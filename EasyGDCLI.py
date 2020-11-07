from GoogleDriverAPI import DrawImage


def cmd_list():
    print('add source')
    print('remove source')
    print('update image')
    print('clear all data')


if __name__ == '__main__':
    di = DrawImage.DrawImage()
    while True:
        cmd_list()
        input_str = input('輸入指令:')
        if input_str.lower() == 'add source':
            di.add_image_source()
        elif input_str.lower() == 'show source':
            di.all_image_source(show=True)
        elif input_str.lower() == 'remove source':
            di.remove_source()
        elif input_str.lower() == 'update image':
            di.update_images()
        elif input_str.lower() == 'clear all data':
            di.clear_data()
        elif input_str.lower() == 'exit':
            exit()
        else:
            print('未知的指令')
