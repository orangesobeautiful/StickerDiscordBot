import discord
import os
from discord.ext import commands
from Database.SQLAlchemyStickerOperation import SQLAlchemyStickerOperation
from GoogleDriverAPI.DrawImage import DrawImage
import pytube
from pytube import YouTube
from apscheduler.schedulers.background import BackgroundScheduler
import random
import time
import pytz


# 頭像提供 https://www.thiswaifudoesnotexist.net/

# BOT_PREFIX = os.environ['prefix'] # -Prefix is need to declare a Command in discord ex: !pizza "!" being the Prefix
# TOKEN = os.environ['token'] # The token is also substituted for security reasons


class StickerBot:
    # sticker_table_name = 'sticker'
    db_operation = None
    draw_image = None
    scheduler = None
    sticker_set = set()
    show_num_max = 3
    my_web_url = os.environ['WebURL']
    # /img-proxy/google-driver
    # image_proxy_url = os.environ['ImgProxyURL']

    cmd_dict = dict()
    cmd_change = False
    bot = None
    bot_prefix = None
    # up is official down is dev
    com_image_em = discord.Embed()
    OPUS_LIBS = ['opus', 'libopus-0.x86.dll', 'libopus-0.x64.dll', 'libopus-0.dll', 'libopus.so.0', 'libopus.0.dylib']
    current_voice_source = None

    def __init__(self):
        self.token = os.environ['bot_token']
        self.db_url = os.environ['DATABASE_URL']
        # self.db_url = "postgres://postgres:@127.0.0.1:5432/postgres"
        self._init_app()

    def start(self):
        self.scheduler.start()
        self.bot.run(self.token)

    def _init_app(self):
        print('discord.opus=' + str(self.load_opus_lib(self.OPUS_LIBS)))
        self._init_db()
        self.scheduler = BackgroundScheduler()
        self.draw_image = DrawImage()
        self.bot_prefix = self.db_operation.get_bot_prefix()
        self.scheduler.add_job(self.all_routine_job, 'cron', hour=5, minute=0, timezone=pytz.timezone('Asia/Taipei'))
        if self.bot_prefix is None:
            self.db_operation.set_bot_prefix('$')
            self.bot_prefix = '$'
        self._init_bot()

    def load_opus_lib(self, opus_libs=OPUS_LIBS):
        # discord.opus
        if discord.opus.is_loaded():
            return True
        else:
            for opus_lib in opus_libs:
                try:
                    discord.opus.load_opus(opus_lib)
                    return True
                except OSError:
                    pass

        raise RuntimeError('Could not load an opus lib. Tried %s' % (', '.join(opus_libs)))

    def _init_db(self):
        self.db_operation = SQLAlchemyStickerOperation(self.db_url)
        self.bot_prefix = self.db_operation.get_bot_prefix()

    def all_routine_job(self):
        self.update_image_warehouse()

    def update_image_warehouse(self):
        print('update image warehouse')
        self.draw_image.update_images()

    """
    def write_sticker_list_file(self):
        with open('./sticker_list.txt', 'w', encoding='utf-8') as sticker_file:
            for cmd in self.cmd_dict:
                sticker_file.write(cmd + '\n')
                sticker_file.write(self.cmd_dict[cmd] + '\n')
    """

    def _init_bot(self):
        self.bot = commands.Bot(command_prefix=self.bot_prefix, description='貼圖小幫手')

        @self.bot.event
        async def on_ready():
            print('discord.__version__ == ' + discord.__version__)
            print('Logged in as ' + self.bot.user.name)
            print('Bot id:' + str(self.bot.user.id))
            print('------')

        @self.bot.event
        async def on_message(msg):
            if msg.author == self.bot.user:  # this is to prevent crashing via infinite loops
                return
            msg_content = msg.content
            msg_channel = msg.channel

            sticker_res = self.db_operation.get_sticker_random(msg_content)
            if sticker_res is not None:
                img_url = sticker_res[0]
                is_gif = sticker_res[1]

                """
                #old method
                if is_gif:
                    await msg_channel.send(img_url)
                else:
                    self.com_image_em.set_image(url=img_url)
                    await msg_channel.send(embed=self.com_image_em)
                """
                await msg_channel.send(img_url)

            await self.bot.process_commands(msg)

        @self.bot.command()
        async def info(ctx):
            embed = discord.Embed(title="貼圖小幫手", description="在樓下幫你支援xxx.jpg\n"
                                                             "使用$help查看指令說明\n"
                                                             "現在支援GoogleDriver共享網址和GIF\n"
                                                             "輕鬆管理貼圖" + self.my_web_url + "\n" +
                                                             "網頁版使用教學" + self.my_web_url + "/sticker-web-tutorial"
                                  , color=0xeee657)

            # 显示机器人所服务的数量。
            # embed.add_field(name="分身數量", value=f"{len(self.bot.guilds)}")
            # 给用户提供一个链接来请求机器人接入他们的服务器
            # embed.add_field(name="Invite", value="[Invite link](<insert your OAuth invitation link here>)")
            await ctx.send(embed=embed)

        self.bot.remove_command('help')

        @self.bot.command(aliases=['help'])
        async def _help(ctx, *args):
            if len(args) == 0:
                embed = discord.Embed(title="貼圖小幫手", description="指令列表:", color=0xeee657)

                embed.add_field(name=self.bot_prefix + "info", value="簡介", inline=False)
                embed.add_field(name=self.bot_prefix + "help", value="指令說明", inline=False)
                embed.add_field(name=self.bot_prefix + "add <貼圖名稱> <圖片網址>", value="將你想要的單字或句子設定成某個圖片網址", inline=False)
                embed.add_field(name=self.bot_prefix + "edit <貼圖ID> <貼圖網址>", value="修改已存在的貼圖", inline=False)
                # embed.add_field(name=self.bot_prefix + "gif <貼圖ID> <設定值>",
                # value="改變貼圖的顯示方式(解決GIF可能顯示不出來的問題)",inline=False)
                embed.add_field(name=self.bot_prefix + "delete <貼圖ID串列>", value="刪除指定ID貼圖", inline=False)
                embed.add_field(name=self.bot_prefix + "deleteST <貼圖名稱>", value="刪除整個貼圖", inline=False)
                embed.add_field(name=self.bot_prefix + "show <貼圖名稱>", value="顯示貼圖的所有圖片資訊", inline=False)
                embed.add_field(name=self.bot_prefix + "exist <貼圖名稱>", value="查詢貼圖是否存在", inline=False)
                embed.add_field(name=self.bot_prefix + "allST", value="顯示所有的貼圖名稱和數量", inline=False)
            else:
                embed = discord.Embed(title="貼圖小幫手", color=0xeee657)
                q_cmd = args[0]
                if q_cmd == 'info':
                    embed.add_field(name=self.bot_prefix + "info", value="簡介", inline=False)
                elif q_cmd == 'help':
                    embed.add_field(name=self.bot_prefix + "help", value="指令說明", inline=False)
                elif q_cmd == 'add':
                    embed.add_field(name=self.bot_prefix + "add <貼圖名稱> <圖片網址>", value="將你想要的單字或句子設定成某個圖片網址\n")
                    embed.add_field(name="<圖片網址>", value="如果是GIF，但不是以.gif結尾\n需要使用 $gif 指令設定其顯示模式", inline=False)
                elif q_cmd == 'edit':
                    embed.add_field(name=self.bot_prefix + "edit <貼圖ID> <貼圖網址>", value="修改已存在的貼圖", inline=False)
                    embed.add_field(name="<貼圖ID>", value="需要整數，可以用 $show 指令查看貼圖ID", inline=False)
                    embed.add_field(name="<圖片網址>", value="如果是GIF，但不是以 .gif 結尾\n需要使用 $gif 指令設定其顯示模式", inline=False)
                elif q_cmd == 'gif':
                    embed.add_field(name=self.bot_prefix + "gif <貼圖ID> <設定值>",
                                    value="改變貼圖的顯示方式(解決GIF可能顯示不出來的問題)\n\n顯示方式分為兩種:\n1.直接顯示(會有網址)\n2.包含在Embed中(不會包含網址列)"
                                          "\n\n如果網址不是以 .gif 結尾，一定需要將其gif值設定為true"
                                    , inline=False)
                    embed.add_field(name="<貼圖ID>", value="需要整數，可以用 $show 指令查看貼圖ID", inline=False)
                    embed.add_field(name="<設定值>", value="(英文大小寫隨意)\n直接顯示: yes, y, true, t, 1"
                                                        "\nEmbed包含: no, n, false, f, 0", inline=False)
                    embed.add_field(name="範例", value="$gif 8763 YeS", inline=False)
                elif q_cmd == 'delete':
                    embed.add_field(name=self.bot_prefix + "delete <貼圖ID串列>", value="刪除指定ID的貼圖", inline=False)
                    embed.add_field(name="<貼圖ID串列>", value="一個或多個整數", inline=False)
                    embed.add_field(name="範例", value="$delete 87\n(刪除多個貼圖)\n$delete 87 63", inline=False)
                elif q_cmd == 'deleteST':
                    embed.add_field(name=self.bot_prefix + "deleteSt <貼圖名稱>", value="刪除整個貼圖", inline=False)
                    embed.add_field(name="<貼圖名稱>", value="就是貼圖的名稱", inline=False)
                elif q_cmd == 'show':
                    embed.add_field(name=self.bot_prefix + "show <貼圖名稱>",
                                    value="顯示貼圖的資訊:\n1.貼圖的圖片數量\n2.貼圖的ID\n3.貼圖的顯示方式", inline=False)
                    embed.add_field(name="<貼圖名稱>", value="就是貼圖的名稱", inline=False)
                elif q_cmd == 'exist':
                    embed.add_field(name=self.bot_prefix + "exist <貼圖名稱>", value="查詢貼圖是否存在", inline=False)
                    embed.add_field(name="<貼圖名稱>", value="就是貼圖的名稱", inline=False)
                elif q_cmd == 'allST':
                    embed.add_field(name=self.bot_prefix + "allST", value="顯示所有的貼圖名稱和數量", inline=False)
                else:
                    await ctx.send('指令 ' + self.bot_prefix + q_cmd + ' 不存在')
                    return

            await ctx.send(embed=embed)

        @self.bot.command()
        async def add(ctx: commands.context.Context, sticker_name: str, img_url: str):
            no_add: list = self.db_operation.add_sticker([{
                'sn': sticker_name,
                'url': img_url,
                'is_gif': False
            }])

            if len(no_add) == 0:
                await ctx.send(sticker_name + ' 新增成功')
            else:
                err_code = no_add[0]['err']
                if err_code == 1:
                    await ctx.send('不支援的網址')
                elif err_code == 2:
                    await ctx.send(sticker_name + ' 已有相同貼圖')

        @self.bot.command()
        async def edit(ctx, sticker_id: str, img_url: str):
            try:
                int(sticker_id)
            except ValueError:
                await ctx.send('需要圖片ID')
                return

            if img_url[-4:] == '.gif':
                is_gif = True
            else:
                is_gif = False

            no_change: list = self.db_operation.edit_sticker([{
                'id': sticker_id,
                'url': img_url,
                'gif': is_gif
            }])

            if len(no_change) == 0:
                await ctx.send('修改成功')
            else:
                err_code = no_change[0]['err']
                if err_code == 1:
                    await ctx.send('不支援的網址')
                elif err_code == 2:
                    await ctx.send('和原本圖片一樣')
                elif err_code == 3:
                    await ctx.send('已有相同圖片')

        @self.bot.command()
        async def gif(ctx, sticker_id: str, is_gif_str: str):
            try:
                int(sticker_id)
            except ValueError:
                await ctx.send("請輸入貼圖ID(數字)")

            is_gif_lower = is_gif_str.lower()
            if is_gif_lower == 't' or is_gif_lower == 'true' or is_gif_lower == '1' or is_gif_lower == 'y' or is_gif_lower == 'yes':
                is_gif = True
            elif is_gif_lower == 'f' or is_gif_lower == 'false' or is_gif_lower == '0' or is_gif_lower == 'n' or is_gif_lower == 'no':
                is_gif = False
            else:
                await ctx.send(is_gif_str + ' 錯誤的格式')
                return

            no_change: list = self.db_operation.edit_sticker([{
                'id': sticker_id,
                'gif': is_gif
            }])

            if len(no_change) == 0:
                await ctx.send('修改成功')
                await ctx.send('注意:新版Discord已經不需要此功能')
            else:
                err_code = no_change[0]['err']
                if err_code == 1:
                    await ctx.send('不支援的網址')
                elif err_code == 2:
                    await ctx.send('和原本圖片一樣')
                elif err_code == 3:
                    await ctx.send('已有相同圖片')

        @self.bot.command()
        async def show(ctx, sticker_name: str):
            img_list = self.db_operation.get_sticker_all(sticker_name)

            if len(img_list) > 0:
                list_len = len(img_list)
                embed = discord.Embed(title=sticker_name, description='總共' + str(list_len) + '個')
                for i, img_ele in enumerate(img_list):
                    img_id = img_ele[0]
                    img_url = img_ele[1]
                    img_gif = img_ele[2]
                    if img_gif:
                        gif_ch_str = '是'
                    else:
                        gif_ch_str = '否'
                    # embed.add_field(name='ID：' + str(img_id), value=img_url + '\ngif:' + gif_ch_str, inline=False)
                    embed.add_field(name='ID：' + str(img_id), value=img_url, inline=False)

                    if i + 1 > self.show_num_max:
                        embed.set_footer(text='還有' + str(list_len - i - 1) + '個未顯示出來')
                        break
                await ctx.send(embed=embed)
            else:
                await ctx.send('貼圖 ' + sticker_name + '不存在')

        @self.bot.command()
        async def delete(ctx, *id_tuple):
            # print(id_tuple)
            abort = False
            if len(id_tuple) > 0:
                id_list = list()
                for x in id_tuple:
                    try:
                        int(x)
                    except ValueError:
                        await ctx.send('請輸入貼圖ID(數字)')
                        abort = True
                        break
                    id_list.append(x)
                if not abort:
                    self.db_operation.delete_sticker(id_list)
                    await ctx.send('執行結束')
            else:
                await ctx.send('請輸入貼圖ID(數字)')

        @self.bot.command()
        async def deleteST(ctx, sticker_name: str):
            self.db_operation.delete_sticker_whole(sticker_name)
            await ctx.send('執行結束')

        @self.bot.command()
        async def exist(ctx, cmd):
            """
            if self.db_operation.get_sticker_first(cmd) is None:
                await ctx.send(cmd + ' 還沒有貼圖')
            else:
                await ctx.send(cmd + ' 已經有了')
            """
            if self.db_operation.is_sticker_name_exist(cmd):
                await ctx.send(cmd + ' 已經有了')
            else:
                await ctx.send(cmd + ' 還沒有貼圖')


        @self.bot.command()
        async def allST(ctx):
            sn_list = self.db_operation.get_all_sn_list()
            send_msg = ''
            send_msg += '貼圖數量' + str(len(sn_list)) + '\n'
            for sticker in sn_list:
                send_msg += sticker + '\t'
            await ctx.send(send_msg)

        @self.bot.command(aliases=['di', '抽卡', '油圖', '抽圖', '抽油圖'])
        async def drawImage(ctx):
            img_id = self.draw_image.get_rand_image_id()
            img_title = 'https://drive.google.com/uc?export=view&id=' + img_id
            """
            #use proxy
            img_url = self.image_proxy_url + '/' + img_id
            """
            img_url = 'https://drive.google.com/uc?export=view&id=' + img_id
            embed = discord.Embed(title=img_title, url=img_url, color=0xff8040)
            embed.set_image(url=img_url)
            if img_url:
                await ctx.send(embed=embed)
            else:
                await ctx.send('還沒有圖片')

        @self.bot.command(pass_context=True)
        async def join(ctx):
            channel = ctx.message.author.voice.channel
            if not channel:
                await ctx.send("You are not connected to a voice channel")
                return
            voice = discord.utils.get(self.bot.voice_clients, guild=ctx.guild)
            if voice and voice.is_connected():
                await voice.move_to(channel)
            else:
                voice = await channel.connect()

        @self.bot.command(aliases=['paly', 'queue', 'que'])
        async def play(ctx, *arg_tuple):
            if len(arg_tuple) == 0:
                await ctx.send('缺少youtube連結')
                return
            else:
                yt_url = arg_tuple[0]

            channel = ctx.message.author.voice.channel
            if not channel:
                await ctx.send("你不在語音頻道中")
                return
            voice = discord.utils.get(self.bot.voice_clients, guild=ctx.guild)
            if voice and voice.is_connected():
                await voice.move_to(channel)
            else:
                voice = await channel.connect()

            try:
                yt = YouTube(yt_url)
                best_url = yt.streams.filter(only_audio=True).order_by(attribute_name='abr')[-1].url
                await ctx.send(':musical_keyboard:`' + yt.title + '`')
            except pytube.exceptions.RegexMatchError:
                await ctx.send("找不到影片")
                return

            guild = ctx.guild
            voice_client: discord.VoiceClient = discord.utils.get(self.bot.voice_clients, guild=guild)
            if voice_client is None:
                await ctx.send("我目前不在語音頻道")
            else:
                if self.current_voice_source is not None:
                    self.current_voice_source.cleanup()

                # self.current_voice_source = await discord.FFmpegOpusAudio.from_probe(best_url)
                self.current_voice_source = await discord.FFmpegOpusAudio.from_probe(best_url)

                if voice_client.is_playing():
                    await ctx.send("已經在播放了")
                else:
                    voice_client.play(self.current_voice_source, after=None)
                    await ctx.send("開始播放")

        @self.bot.command()
        async def stop(ctx):
            guild = ctx.guild
            voice_client: discord.VoiceClient = discord.utils.get(self.bot.voice_clients, guild=guild)
            if voice_client.is_playing():
                voice_client.stop()
                self.current_voice_source.cleanup()
                await ctx.send("已經停止")
            else:
                await ctx.send("目前沒有音樂")

        @self.bot.command(aliases=['exit'])
        async def disconnect(ctx):
            guild = ctx.guild
            voice_client: discord.VoiceClient = discord.utils.get(self.bot.voice_clients, guild=guild)
            try:
                await voice_client.disconnect()
            except AttributeError:
                pass


if __name__ == '__main__':
    stb = StickerBot()
    stb.start()
