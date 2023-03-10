from dotenv import load_dotenv
from discord.ext import commands
import discord
import datetime
import aiohttp
import os

# setting initial vars
load_dotenv()
auth_header = os.getenv('REQUESTS_AUTH')
token = os.getenv('DISCORD_TOKEN')
url = os.getenv('DEFAULT_DISC_URL')
intents = discord.Intents.all()
bot = commands.Bot(command_prefix='!', intents=intents)
ctx = commands.Context
client = discord.Client(intents=intents)

# testing syn & ack for bot comms
@bot.command()
async def hello(ctx):
    await ctx.send('hello and welcome')

# eventual command to get user info and loop through each user to
# see when last activity was
@bot.command()
async def users(ctx):
    member_list = []

    for guild in bot.guilds:
        for member in guild.members:
            if member.bot == True:
                continue
            else:
                member_list.append(member.name)

    print(member_list)
    await ctx.send(member_list)

# testing retrieving all channel IDs
@bot.command()
async def channels(ctx):
    channel_list = []

    for guild in bot.guilds:
        for channel in guild.channels:
            if channel.name == "Text Channels":
                continue
            elif channel.name == "Voice Channels":
                continue
            else:
                channel_list.append(channel)

    print(channel_list)
    await ctx.send(channel_list)

# getting message history, playing around with datetimes and deltas
@bot.command()
async def test(ctx):
    headers = {'authorization': auth_header,
               'content_type': 'application/json'}
    url = "https://discord.com/api/v9/channels/1081830805288001576/messages"
    async with aiohttp.ClientSession(headers=headers) as session:
        async with session.get(url=url) as req:
            response = await req.json()
            print(response[0]['author'])
            await ctx.send(response[0]['author'])

    # message_timestamps = []
    # for message in messages:
    #     message_timestamps.append(message["timestamp"])

    #     old_timestamp = message_timestamps[0]
    #     a = old_timestamp.split('+')
    #     b = a[0]
    #     c = b.replace('T', ' ')
    #     newtime = datetime.datetime.strptime(c, "%Y-%m-%d %H:%M:%S.%f")
    #     endtime = newtime + datetime.timedelta(days=-10)
    #     print(endtime)
    # #print(datetime.now())

#asyncio.run(retrieve_messages('1081830736480436295'))

# event that checks when someone has joined a voice channel
@bot.event
async def on_voice_state_update(member, before, after):
    current_date_utc = datetime.datetime.now()
    voice_message_on_join = f'{member} joined voice channel {after.channel} at {str(current_date_utc)}'
    
    if before.channel is None and after.channel is not None:
        get_log_channel = await bot.fetch_channel('1083587191336341654')
        await get_log_channel.send(voice_message_on_join)

bot.run(token)