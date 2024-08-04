import { Button, FormControl, MenuItem, TextField } from "@mui/material";

function QueryView(props) {
  const from = props.from;
  const setFrom = props.setFrom;
  const filter = props.filter;
  const setFilter = props.setFilter;
  const groupBy = props.groupBy;
  const setGroupBy = props.setGroupBy;

  return (
    <div>
      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField
          label="Kind"
          select
          value={props.queryKind}
          size="small"
          onChange={(event) => {
            props.setQueryKind(event.target.value)
          }}
          InputLabelProps={{ shrink: true }}
        >
          {
            props.kinds.map((value) => (
              <MenuItem value={value}>{value}</MenuItem>
            ))
          }
        </TextField>
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="From" value={from} onChange={(event) => {
          setFrom(event.target.value)
        }} InputLabelProps={{ shrink: true }} />
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="Filter" value={filter} onChange={(event) => {
          setFilter(event.target.value)
        }} InputLabelProps={{ shrink: true }} />
      </FormControl>

      <FormControl fullWidth sx={{marginBottom: '10px'}}>
        <TextField size="small" label="Group By" value={groupBy} onChange={(event) => {
          setGroupBy(event.target.value)
        }} InputLabelProps={{ shrink: true }} />
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
