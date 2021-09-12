import os
import sys
import argparse
import flask
import flask_login
import json
from flask import Flask, request, abort, render_template, jsonify, send_from_directory, session
from flask_sslify import SSLify
from flask_cors import CORS
from Database.SQLAlchemyStickerOperation import SQLAlchemyStickerOperation
from Database.SQLAlchemyWebLoginOperation import SQLAlchemyWebLoginOperation
from CommonFunction import StickerCommon

db_url = ''
secret_key = ''
router_prefix = ''
image_return = False
save_image_local = False
sticker_download_dir = StickerCommon.sticker_dir


def _read_setting():
    global db_url, image_return, save_image_local, secret_key, router_prefix
    _, db_url, _, _, save_image_local, image_return, _, secret_key, router_prefix = StickerCommon.read_setting()
    if db_url == '' or save_image_local is None or image_return is None:
        return False
    return True


if not _read_setting():
    raise OSError('讀取設定失敗，需要在環境變數或設定檔(setting.ini)中提供完整設定值')
sticker_db_operation = SQLAlchemyStickerOperation(db_url, save_image_local)
web_login_db_operation = SQLAlchemyWebLoginOperation(db_url)

app = Flask(__name__)
app.secret_key = secret_key
app.config.from_object(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})

# sslify = SSLify(app, skips=['h'])
# SECURE_REDIRECT_EXEMPT = ['/']

login_manager = flask_login.LoginManager()
login_manager.login_message_category = 'info'
login_manager.login_message = 'Access denied.'
login_manager.init_app(app)

"""
@app.after_request
def after_request(resp):
    resp = Flask.make_response(resp)
    resp.headers['Access-Control-Allow-Origin'] = 'http://localhost:5000'
    resp.headers['Access-Control-Allow-Methods'] = 'GET,POST'
    resp.headers['Access-Control-Allow-Headers'] = 'content-type,token'
    return resp
"""


class User(flask_login.UserMixin):
    pass


@login_manager.user_loader
def load_user(user_id):
    user = User()
    user.id = user_id
    return user


if image_return:
    @app.route("/sticker-image/<path:filename>", methods=['GET'])
    def sticker_image_return(filename):
        img_path = os.path.join(sticker_download_dir, filename)
        if os.path.isfile(img_path):
            with open(img_path, 'rb') as img_file:
                img_b = img_file.read()

            last_dot_index = filename.rfind('.')
            if last_dot_index >= 0:
                img_type = filename[last_dot_index + 1:].lower()
            else:
                img_type = 'jpeg'
            return img_b, 200, {'content-type': 'image/' + img_type,
                                'Content-Disposition': 'inline;filename="' + filename + '";filename*=UTF-8\'\'' + filename}
        else:
            return '404', 404


@app.route(router_prefix + '/has_login', methods=['GET'])
@flask_login.login_required
def has_login():
    return '', 200


@app.route(router_prefix + '/get_login_code', methods=['GET'])
def get_login_code():
    login_code, expiration_time = web_login_db_operation.generate_verification_code()
    return jsonify({'code': login_code})


@app.route(router_prefix + '/check_login', methods=['GET'])
def check_login():
    try:
        code = request.args.get('code')
    except ValueError:
        return jsonify({'error': '錯誤的參數'})
    login_success, user_id = web_login_db_operation.check_verification_status(code)
    if login_success:
        user = User()
        user.id = user_id
        flask_login.login_user(user)
        web_login_db_operation.delete_code(code)
        return jsonify({'result': '1', 'user_id': user_id})
    else:
        return jsonify({'result': '0', 'user_id': -1})


@app.route(router_prefix + '/user_info', methods=['GET'])
@flask_login.login_required
def user_info():
    return_data = web_login_db_operation.get_user_info(flask_login.current_user.id)
    if return_data is None:
        return '404', 404
    else:
        name = return_data[0]
        avatar_url = return_data[1]
    return jsonify({'name': name, 'avatar_url': avatar_url})


@app.route(router_prefix + '/logout')
def logout():
    flask_login.logout_user()
    return 'Logged out'


@app.route(router_prefix + "/all_sticker", methods=['GET'])
@flask_login.login_required
def all_sticker():
    try:
        start = request.args.get('start')
        num = request.args.get('num')
        int(start)
        int(num)
    except ValueError:
        return jsonify({'error': '錯誤的參數'})

    # print('start', start)
    # print('start', start)
    # print('num', num)
    get_res_list = sticker_db_operation.get_sticker_group_by_name(start=start, num=num)
    maxp = sticker_db_operation.max_page(num)
    return_json = dict()
    # print(maxp)
    return_json['maxp'] = maxp
    return_json['img_data'] = list()
    # print(get_res_list)
    for sticker_ele in get_res_list:
        sticker_name = sticker_ele[0]
        if type(sticker_ele[1]) == str:
            sticker_list = json.loads(sticker_ele[1])
        else:
            sticker_list = sticker_ele[1]
        for sticker_prop in sticker_list:
            is_gif = sticker_prop['gif']
            if type(is_gif) == int:
                sticker_prop['gif'] = bool(sticker_prop['gif'])
        return_json['img_data'].append({'sn': sticker_name, 'sts': sticker_list})

    return jsonify(return_json)


@app.route(router_prefix + "/search", methods=['GET'])
@flask_login.login_required
def search():
    try:
        query = request.args.get('q')
    except ValueError:
        return jsonify({'error': '錯誤的參數'})

    get_res_list = sticker_db_operation.get_sticker_all(query)

    return_json = dict()
    sn_data_dict_list = list()
    for sticker_ele in get_res_list:
        sn_data_dict_list.append({'id': sticker_ele[0], 'url': sticker_ele[1], 'gif': sticker_ele[2]})

    return_json['maxp'] = 1
    return_json['img_data'] = list()

    if len(get_res_list) >= 1:
        return_json['img_data'].append(
            {'sn': query, 'sts': sn_data_dict_list})

    return jsonify(return_json)


@app.route(router_prefix + "/change_sn", methods=['POST'])
@flask_login.login_required
def change_sn():
    if request.method == 'POST':
        change_request = request.json
        sticker_name = change_request['sn']
        # print(change_request)
        for rtype in change_request:
            if rtype == 'add' and len(change_request[rtype]) > 0:
                add_list = list()
                for add_img in change_request[rtype]:
                    add_list.append({'sn': sticker_name, 'url': add_img['url'], 'is_gif': add_img['gif']})
                sticker_db_operation.add_sticker(add_list)
            elif rtype == 'edit' and len(change_request[rtype]) > 0:
                edit_list = list()
                for edit_img in change_request[rtype]:
                    edit_list.append(edit_img)
                sticker_db_operation.edit_sticker(edit_list)
            elif rtype == 'delete' and len(change_request[rtype]) > 0:
                del_list = list()
                for edit_img in change_request[rtype]:
                    del_list .append(edit_img)
                sticker_db_operation.delete_sticker(del_list)

        r_data = dict()
        # 把 tuple 轉為 list ， jsonify 不支援 tuple
        img_list = list()
        for sticker_ele in sticker_db_operation.get_sticker_all(sticker_name):
            img_list.append(list(sticker_ele))
        r_data['imgs'] = img_list
        r_data['err'] = ''
        return jsonify(r_data)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default='127.0.0.1', help='Listen Host')
    parser.add_argument("--port", default='5000', help='Listen Port')
    parser.add_argument("-proxy", default=False, help='Has Proxy Header?', action="store_true")
    parser.add_argument("-debug", default=False, help='debug mode?', action="store_true")
    args = parser.parse_args()

    if args.proxy:
        has_proxy = args.proxy

    print("Flask Version:" + flask.__version__)
    app.debug = args.debug
    app.run(host=args.host, port=int(args.port))
    print('next')