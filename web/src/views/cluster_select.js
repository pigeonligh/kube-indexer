import * as React from 'react';
import { FormControl, MenuItem, TextField } from "@mui/material";
import { ChangeCluster, GetCluster } from "../js/cluster";

function ClusterSelect(props) {
  const [clusters, setClusters] = React.useState([]);

  React.useEffect(() => {
    fetch('/api/cluster', {
      method: 'GET'
    }).then(
      (response) => response.json()
    ).then((data) => {
      setClusters(data)
    }).catch((err) => {
      console.log(err.message);
    });
  }, [props])

  return (
    <FormControl fullWidth>
      <TextField
        label="Cluster"
        select
        value={GetCluster()}
        size="small"
        onChange={(event) => {
          ChangeCluster(event.target.value)
        }}
        InputLabelProps={{ shrink: true }}
      >
        {
          clusters.map((value) => (
            <MenuItem value={value}>{value}</MenuItem>
          ))
        }
      </TextField>
    </FormControl>
  );
}

export default ClusterSelect;
