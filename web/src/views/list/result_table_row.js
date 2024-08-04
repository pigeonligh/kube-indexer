import * as React from 'react';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import ResultObjectViewer from './result_object_viewer';
import { Button, Modal } from '@mui/material';
import RawViewer from './raw_viewer';

function ResultTableRow(props) {
  const key = props.data.key;
  const row = props.data.value;
  const headers = props.headers;
  const [collapseOpen, setCollapseOpen] = React.useState(false);
  const [modalOpen, setModalOpen] = React.useState(false);

  const getTableRow = () => {
    const arr = [];
    arr.push(
      <TableCell>
        <IconButton
          aria-label="expand row"
          size="small"
          onClick={() => setCollapseOpen(!collapseOpen)}
        >
          {collapseOpen ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
        </IconButton>
      </TableCell>
    );
    headers.forEach(element => {
      arr.push(
        <TableCell key={element}> {row[element]} </TableCell>
      );
    });
    arr.push(
      <TableCell key="_raw">
        <Button data-nocollapse={true} onClick={() => {setModalOpen(true)}}>raw</Button>
      </TableCell>
    )
    return arr
  }

  return (
    <React.Fragment>
      <TableRow
        key={ key }
        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
        onClick={(event) => {
          if (!event.target.dataset.hasOwnProperty('nocollapse')) {
            setCollapseOpen(!collapseOpen)
          }
        }}
      >
        {getTableRow()}
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0, background: 'white' }} colSpan={6}>
          <Collapse in={collapseOpen} timeout="auto" unmountOnExit>
            <ResultObjectViewer data={row} queryFn={props.queryFn}/>
          </Collapse>
        </TableCell>
      </TableRow>
      <Modal
        open={modalOpen}
        onClose={() => {setModalOpen(false)}}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <RawViewer kind={row._kind} resourcekey={key} open={modalOpen}/>
      </Modal>
    </React.Fragment>
  );
}

export default ResultTableRow;
