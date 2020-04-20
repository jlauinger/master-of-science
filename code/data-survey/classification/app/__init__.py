from flask import Flask
from app.config import Config
from flask_bootstrap import Bootstrap

app = Flask(__name__)
app.config.from_object(Config)

Bootstrap(app)

from app import routes