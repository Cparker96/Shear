import logging

logger = logging.getLogger()

# LIST OF LOGGING LEVELS WITH STR AND INT REPRESENTATIONS IN ORDER:

# LEVEL   || NUMERIC VALUE
# ========================
# NOTSET  ||     0
# DEBUG   ||    10
# INFO    ||    20
# WARNING ||    30
# ERROR   ||    40
# CRITIAL ||    50

# creating handlers
#declare_stream_handler = logging.StreamHandler()
info_handler = logging.FileHandler('valid_logs.log')
error_handler_voice = logging.FileHandler('invalid_voice_channel_logs.log')
#declare_stream_handler.setLevel('DEBUG')
info_handler.setLevel('INFO')
error_handler_voice.setLevel('ERROR')

# creating formatters
#declare_stream_format = logging.Formatter('%(asctime)s - %(name)s - %(message)s')
info_format = logging.Formatter('%(asctime)s - %(levelname)s - %(name)s - %(message)s')
error_format = logging.Formatter('%(asctime)s - %(levelname)s - %(name)s - %(message)s')

# set the formatters
#declare_stream_handler.setFormatter(declare_stream_format)
info_handler.setFormatter(info_format)
error_handler_voice.setFormatter(error_format)

# add handlers to the logger
#logger.addHandler(declare_stream_handler)
logger.addHandler(info_handler)
logger.addHandler(error_handler_voice)


