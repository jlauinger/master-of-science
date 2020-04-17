#!/usr/bin/env python

from .analysis_steps import *


offset = int(sys.argv[1])
length = int(sys.argv[2])

projects_df = get_projects_data(offset, length)

for i, project in projects_df.iterrows():
    try:
        if not go_mod_exists(project):
            print(str(i+1-offset) + "/" + str(length) + " (#" + str(i+1) + "): Skipping " + project['project_name'])
            ErrorCondition(stage='go.mod',
                           project_name=project['project_name'],
                           message='go.mod not found',
                           module_import_path='',
                           file_name='').write()
            continue

        print(str(i+1-offset) + "/" + str(length) + " (#" + str(i+1) + "): Analyzing " + project['project_name'])

        analyze_modules(project)

    except Exception as e:
        ErrorCondition(stage='project',
                       project_name=project['project_name'],
                       message="{}".format(e),
                       module_import_path='',
                       file_name='').write()
        print('SAVING AN ERROR!')
        continue

# if we actually ever come to this, close the open files
Module.close()
MatchResult.close()
VetResult.close()
ErrorCondition.close()
