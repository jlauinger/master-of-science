from flask_wtf import FlaskForm
from wtforms import StringField, SubmitField, TextAreaField

class ClassificationForm1(FlaskForm):
    label1 = StringField('Label 1 (What)')
    submit1 = SubmitField('Classify!')

class ClassificationForm2(FlaskForm):
    label2 = StringField('Label 2 (Purpose)')
    submit2 = SubmitField('Classify!')