# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import shutil
import logging
from sys import implementation, meta_path
from mkdocs import plugins
import yaml
import os
import pandas
from fnmatch import fnmatch
import glob

log = logging.getLogger('mkdocs')

@plugins.event_priority(100)
def on_pre_build(config, **kwargs):
    log.info("copying geps")
    shutil.copytree("geps","site-src/geps", dirs_exist_ok=True)

    # calling to get the conformance reports generated
    yamlReports = getYaml()
    
    create_md(yamlReports)

# outputs reports to markdown file
def create_md(reports):
    # getting rid of some columns
    reports = reports.drop(columns=['implementation']) 
    
    tests = reports[["organization","version","name", "extended.supportedFeatures"]]
    tests = tests.groupby(['organization']).apply(lambda x: x)

    testNames = tests['name'].unique() # HTTP, TLS, MESH, etc.
    df =tests

    table= reports.groupby(["organization"], as_index=False).name.apply(' '.join).apply(lambda x: x)

    for n in testNames:
        temp = df.loc[df['name']==n]
        temp.rename(columns={"extended.supportedFeatures":n+': Supported Features'},inplace=True)
        temp=temp.drop(["name","organization"],axis=1)
        temp.reset_index(inplace=True)
        temp = temp.drop(["level_1"],axis=1)
        table = table.merge(temp, how="left")

    # dropping TLS supportedFeatures column since no implementation has listed any supported features
    table = table.drop(["TLS: Supported Features"], axis=1)
    table.rename(columns={"organization":"Organization", "name":"Protocol Profile","version":"Version" }, inplace=True)
    table = table.fillna("N/A")
    
    # keep the latest version in the table
    table.sort_values(['Organization','Version'], inplace=True)
    table.drop_duplicates(subset="Organization", inplace=True,keep='last')

    # Output markdown table
    with open('site-src/implementation-table.md','w') as f:
        f.write("This table is populated from the conformance reports uploaded by project implementations.\n\n")
        f.write(table.to_markdown(index=False)+'\n')



# the path should be changed when there is a new version
conformance_path = "conformance/reports/v1.0.0/**"
def getYaml():
    log.info("parsing conformance reports ============================")
    yamls = []

    # reports must be named according to the following pattern : <API Channel>-<Implementation version>-<mode>-report.yaml

    for p in glob.glob(conformance_path, recursive=True): # getting all the paths in conforamnce

        if fnmatch(p, "*.yaml"):
            
            x = load_yaml(p)
            profiles = pandas.json_normalize(x, record_path='profiles',meta=["implementation"] ) 
            
            implementation = pandas.json_normalize(profiles.implementation)
            yamls.append(pandas.concat([implementation,profiles], axis=1))

    yamls = pandas.concat(yamls)
    return yamls

def load_yaml(name):
    with open(name, 'r') as file:
        x = yaml.safe_load(file)

    return x

