import os

class Config(object):
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'you-will-never-guess'
    DATA_DIR = os.environ.get('DATA_PATH') or '/home/johannes/studium/s14/masterarbeit/code/data-survey/data'
    GO_MOD_PATH = os.environ.get('GO_MOD_PATH') or '/home/johannes/.go'
    GO_LIB_PATH = os.environ.get('GO_LIB_PATH') or '/usr/lib/go'