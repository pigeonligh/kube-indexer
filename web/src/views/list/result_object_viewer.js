import { Chip, Link, List, ListItem, ListItemButton, Table, TableBody, TableCell, TableContainer, TableRow } from '@mui/material';
import JsonView from '@uiw/react-json-view';
import { githubLightTheme } from '@uiw/react-json-view/githubLight';

function ResultObjectViewer(props) {
  const data = JSON.parse(JSON.stringify(props.data));

  delete data._raw;
  delete data._key;
  delete data._kind;
  delete data._resource_version;

  const viewListItem = (k, val) => {
    if (val === null) {
      return (<Chip size="small" label="NULL"/>)
    }
    if (typeof val === 'object') {
      if (val.hasOwnProperty('ref')) {
        return (
          <li>
            <Link onClick={() => {
              props.queryFn(
                val.ref.kind,
                "global."+props.data._kind+"[\""+props.data._key+"\"]."+k,
                "",
                "",
              )
            }}>
              {val.ref.kind}
            </Link> / <Link onClick={() => {
              props.queryFn(
                val.ref.kind,
                "global."+props.data._kind+"[\""+props.data._key+"\"]."+k,
                "cur._key==\""+val.ref.key+"\"",
                "",
              )
            }}>
              {val.ref.key}
            </Link>
          </li>
        )
      }

      return (
        <JsonView 
          value={val}
          style={githubLightTheme}
          displayDataTypes={false}
          collapsed={false}
          shortenTextAfterLength={0}
          objectSortKeys={true}
        />
      )
    }
    return (<span>{val}</span>)
  }

  const viewValue = (k, val) => {
    if (val === null) {
      return (<Chip size="small" label="NULL"/>)
    }
    if (typeof val === 'object') {
      if (Array.isArray(val)) {
        if (val.length > 0) {
          return (
            <List>
              {val.map((v) => viewListItem(k, v))}
            </List>
          )
        }
        return (<Chip size="small" label="NONE"/>)
      }
      return (
        <JsonView 
          value={val}
          style={githubLightTheme}
          displayDataTypes={false}
          collapsed={false}
          shortenTextAfterLength={0}
          objectSortKeys={true}
        />
      )
    }
    return (<span>{val}</span>)
  }

  const getTableBody = () => {
    if (data) {
      const keys = Object.keys(data);
      console.log(Object.keys(data))
      return (
        <TableBody>
          {keys.map((k) => (
            <TableRow>
              <TableCell>
                {k}
              </TableCell>
              <TableCell>
                {viewValue(k, data[k])}
              </TableCell>
          </TableRow>
          ))}
        </TableBody>
      )
    }
  }

  return (
    <div>
      <TableContainer>
        <Table size="small">
          {getTableBody()}
        </Table>
      </TableContainer>
    </div>
  );
}

export default ResultObjectViewer;
