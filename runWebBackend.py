from web.backend import StickerManager
from bot import BotStart
import sys


if __name__ == '__main__':
    StickerManager.app.debug = False
    StickerManager.app.run("127.0.0.1")
