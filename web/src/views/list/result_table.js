import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import ResultTableRow from './result_table_row';

function ResultTable(props) {
  const headers = props.headers;

  const getTableHeadRows = () => {
    const arr = [];
    arr.push(
      <TableCell />
    )
    headers.forEach(element => {
      arr.push(
        <TableCell key={element}> {element} </TableCell>
      );
    });
    arr.push(
      <TableCell />
    )
    return (
      <TableRow>
        {arr}
      </TableRow>
    )
  }

  const getTableData = (row) => {
    return <ResultTableRow headers={headers} data={row}/>
  }

  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          {getTableHeadRows(props)}
        </TableHead>
        <TableBody>
          {props.data.map(getTableData)}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default ResultTable;
