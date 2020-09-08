# Go snippet classification tool

This is a Python Flask application that I used to classify Go code snippets. It modifies CSV data at its core, but
provides a nice user interface to me.


## Deployment

 1. Create a virtual environment for the app: `python3 -m venv venv`
 2. Activate the environment: `source venv/bin/activate`
 3. Install dependencies: `pip install -r requirements.txt`
 4. Configure the app by editing the paths contained in `classification.ini`
 5. Install [uwsgi](https://uwsgi-docs.readthedocs.io/en/latest/WSGIquickstart.html)
 6. Create a Systemd service or something similar to run the uwsgi configuration provided by `classification.ini`
 7. Configure a webserver, e.g. nginx, to forward traffic at some URL to the uwsgi socket for the app
 8. Visit the URL to use the application
 

## Usage

This application heavily uses internal state because that was absolutely acceptable for the purpose it was built for,
i.e. I was the only user. Most importantly, links to snippets contain an index, but this index is relative to the
current data file. If the file is switched, all links now also point to a different snippet. Therefore, links can only
be shared referring to the same snippet for short times, they will not persist. Also, if the service is restarted before
data was saved, all modifications are lost.

First, select a data file to use with the button at the top of the index page. You can also switch whether you are
classifying the first or second dimension there, this will affect which snippets are in the todo section and which are
shown as finished.

Then, click on a link to a snippet to view the snippet classification page. It shows meta information about the
snippet as well as the code and context.

Enter the label you want to save for this snippet in the text box below the code and hit the classify button to assign
the label and go on to the next snippet. To assign both labels you have to go back after the first and assign the second
individually. For every label that already exists for some snippet in the currently selected data file, there is a
quick-assign button.

**Note**: due to efficiency reasons, label modifications do not actually get saved to disk until you click the Save
button on the snippets index list.


## Code structure

The application code lives within the `app/` directory, the files in the root directory are only Flask boilerplate code.

 - `app/config.py` loads configuration options from the environment (e.g. the uwsgi `classification.ini` file) or sets
   default values
 - `app/data_io.py` is responsible for loading and saving CSV data to disk using Pandas
 - `app/forms.py` contains Flask form definitions for the label assignment forms
 - `app/routes.py` contains controller code for all GET/POST routes. This is most of the actual application logic
 - `app/templates/` contains Jinja2 template code to generate the HTML
