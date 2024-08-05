import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import ResultTableRow from './result_table_row';
import { TablePagination } from '@mui/material';

function ResultTable(props) {
  const headers = props.headers;
  const [rowsPerPage, setRowsPerPage] = React.useState(10);
  const [page, setPage] = React.useState(0);

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
    return <ResultTableRow headers={headers} data={row} queryFn={props.queryFn}/>
  }

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  return (
    <div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650, background: '#f7f7f7' }}>
          <TableHead>
            {getTableHeadRows(props)}
          </TableHead>
          <TableBody>
            {props.data.slice(
              page * rowsPerPage,
              page * rowsPerPage + rowsPerPage,
            ).map(getTableData)}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 20, 40]}
        component="div"
        count={props.data.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
      />
    </div>
  );
}

export default ResultTable;
