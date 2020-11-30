import os
import configparser

project_dir = os.path.abspath(os.path.join(os.path.dirname(os.path.realpath(__file__)), os.path.pardir))
sticker_dir = os.path.join(project_dir, 'sticker-image')
if not os.path.isdir(sticker_dir):
    os.mkdir(sticker_dir)


def str_to_bool(s: str):
    if s.lower() == 'true':
        return True
    elif s.lower() == 'false':
        return False
    else:
        return None


def read_setting():
    token = ''
    db_url = ''
    my_web_url = ''
    sticker_url = ''
    save_image_local = None
    flask_return_sticker = None

    read_success = False
    try:
        token = os.environ['bot_token']
        db_url = os.environ['DATABASE_URL']
        my_web_url = os.environ['WebURL']
        sticker_url = os.environ['StickerURL']
        save_image_local = os.environ['SaveImageLocal']
        flask_return_sticker = os.environ['FlaskReturnSticker']
        save_image_local = str_to_bool(save_image_local)
        flask_return_sticker = str_to_bool(flask_return_sticker)
        read_success = True
    except KeyError:
        pass

    if not read_success:
        setting_path = os.path.join(project_dir, 'setting.ini')
        conf = configparser.ConfigParser()
        if not os.path.isfile(setting_path):
            conf.add_section('Environment')
            conf.set('Environment', 'bot_token', '')
            conf.set('Environment', 'DATABASE_URL', '')
            conf.set('Environment', 'WebURL', '')
            conf.set('Environment', 'StickerURL', '')
            conf.add_section('AdditionFunction')
            conf.set('AdditionFunction', 'SaveImageLocal', 'False')
            conf.set('AdditionFunction', 'FlaskReturnSticker', 'False')
            with open(setting_path, 'w', encoding='utf-8') as setting_file:
                conf.write(setting_file)
        else:
            try:
                conf.read(setting_path, encoding='utf-8')
                env_section = conf['Environment']
                token = env_section['bot_token']
                db_url = env_section['DATABASE_URL']
                my_web_url = env_section['WebURL']
                sticker_url = env_section['StickerURL']
                addition_section = conf['AdditionFunction']
                save_image_local = addition_section['SaveImageLocal']
                flask_return_sticker = addition_section['FlaskReturnSticker']
                save_image_local = str_to_bool(save_image_local)
                flask_return_sticker = str_to_bool(flask_return_sticker)
            except KeyError:
                pass

    return token, db_url, my_web_url, sticker_url, save_image_local, flask_return_sticker
