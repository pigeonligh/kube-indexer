import * as React from 'react';
import { Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";

function QueryView(props) {
  const [from, setFrom] = React.useState("");
  const [filter, setFilter] = React.useState("");
  const [groupBy, setGroupBy] = React.useState("");

  return (
    <div>
      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <InputLabel id="demo-simple-select-label">Kind</InputLabel>
        <Select
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={props.queryKind}
          label="Kind"
          size="small"
          onChange={(event) => {
            props.setQueryKind(event.target.value)
          }}
        >
          {
            props.kinds.map((value) => (
              <MenuItem value={value}>{value}</MenuItem>
            ))
          }
        </Select>
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="From" value={from} onChange={(event) => {
          setFrom(event.target.value)
        }} />
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="Filter" value={filter} onChange={(event) => {
          setFilter(event.target.value)
        }} />
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="Group By" value={groupBy} onChange={(event) => {
          setGroupBy(event.target.value)
        }} />
      </FormControl>

      <FormControl fullWidth>
        <Button onClick={() => {
          props.queryFn(props.queryKind, from, filter, groupBy)
        }}>Search</Button>
      </FormControl>
    </div>
  );
}

export default QueryView;
