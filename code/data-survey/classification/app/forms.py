from flask_wtf import FlaskForm
from wtforms import StringField, SubmitField, TextAreaField

class ClassificationForm(FlaskForm):
    label = StringField('Label')

    submit = SubmitField('Classify!')