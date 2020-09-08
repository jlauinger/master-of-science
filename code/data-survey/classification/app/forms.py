from flask_wtf import FlaskForm
from wtforms import StringField, SubmitField, TextAreaField

# form to save label
class ClassificationForm1(FlaskForm):
    label1 = StringField('Label 1 (What)')
    submit1 = SubmitField('Classify!')

# form to save label2
class ClassificationForm2(FlaskForm):
    label2 = StringField('Label 2 (Purpose)')
    submit2 = SubmitField('Classify!')