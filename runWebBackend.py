from web.backend import web_test
from bot import BotStart
import sys


if __name__ == '__main__':
    web_test.app.debug = False
    web_test.app.run("127.0.0.1")
