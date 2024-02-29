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
from mkdocs import plugins
import yaml
import os
import pandas

log = logging.getLogger('mkdocs')

@plugins.event_priority(100)
def on_pre_build(config, **kwargs):
    log.info("copying geps")
    shutil.copytree("geps","site-src/geps", dirs_exist_ok=True)

    # calling to get the conformance reports generated
    reports = generate()
    create_md(reports)

def create_md(reports):
    pandas.set_option('display.max_colwidth',0)

    data = pandas.DataFrame.from_dict(reports[0])
    with open('site-src/test.md','w') as f:
        f.write(data.to_markdown())




# TODO : versioning of the reports should be changed
conformance_path = "conformance/reports/v1.0.0/"
def generate():
    log.info("parsing conformance reports ============================")
    yamls = []
    for f in os.listdir(conformance_path):
        x = load_yaml(conformance_path+f)
        yamls.append(pandas.json_normalize(x))
    return yamls

def load_yaml(name):
    with open(name, 'r') as file:
        x = yaml.safe_load(file)

    return x

