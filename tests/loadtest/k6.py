#
# * Copyright (c) 2021, arivum.
# * All rights reserved.
# * SPDX-License-Identifier: MIT
# * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
#

#%%
import json
import pandas
import altair
import altair_saver
from datetime import datetime

from pandas.core.reshape.reshape import stack

seconds = 20

def read_and_transform(jsonFile):
    data = []
    with open(jsonFile, 'r') as k6raw:
        for line in k6raw:
            j = json.loads(line.rstrip())
            if j["type"] == "Point" and j["metric"] == "http_req_failed" and j["data"]:
                try:
                    ts = datetime.strptime(j["data"]["time"].split(".")[0], '%Y-%m-%dT%H:%M:%S')
                    data.append({"date": ts, "value": j["data"]["value"], "status": ("succeeded" if j["data"]["value"] == 0 else "failed")})
                except:
                    continue
    data = pandas.DataFrame(data)
    data = data.groupby(['date', 'status']).count().reset_index().groupby('status').resample(str(seconds)+"s", on="date").sum().reset_index()
    data['date'] = (data['date'] - data['date'].min()).apply(lambda x: x.total_seconds())
    data['value'] = data['value'].apply(lambda x: x/seconds)
    return data




# %%
altair.data_transformers.disable_max_rows()

def draw(data, title):
    return altair.Chart(data, title=title).mark_area(
        line=True,
    ).encode(
        x=altair.X('date:Q', title="Time [s]"),
        y=altair.Y('value:Q', title="Rate [req/s]",stack="zero", scale=altair.Scale(domain=(0, 25))),
        color=altair.Color('status:N', scale=altair.Scale(scheme='darkmulti')),
    )
# failed_req = altair.Chart(data).mark_area(
#     color="red",
#     line=True
# ).encode(
#     x=altair.X('date:T'),
#     y=altair.Y('sum(value):Q', title="req/s"),
# )

# %%
filename_base="base"
chart = draw(read_and_transform("results/"+filename_base+".json"), "Throughput without dynratelimiter") | draw(read_and_transform("results/"+filename_base+"_rate_limited.json"), "Throughput with dynratelimiter")
chart
altair_saver.save(chart, filename_base+".png")
#+ failed_req
# %%
