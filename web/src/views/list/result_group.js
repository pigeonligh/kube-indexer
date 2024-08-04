import * as React from 'react';
import { Box, Button, Collapse } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import ResultTable from './result_table'

function ResultGroup(props) {
  const [collapseOpen, setCollapseOpen] = React.useState(false);

  const getGroupView = () => {
    if (props.data.name === "-") {
      return (
        <ResultTable 
          name={props.data.name}
          headers={props.headers}
          data={props.data.items}
        />
      )
    }
    return (
      <Box>
        <IconButton
          aria-label="expand row"
          size="large"
          onClick={() => setCollapseOpen(!collapseOpen)}
        >
          {collapseOpen ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
        </IconButton>
        <Button
          size="large"
          sx={{textTransform: 'none'}}
          onClick={() => setCollapseOpen(!collapseOpen)}
        >
          {props.data.name} ({props.data.count})
        </Button>
        

        <Collapse in={collapseOpen} timeout="auto" unmountOnExit>
          <ResultTable 
            name={props.data.name}
            headers={props.headers}
            data={props.data.items}
          />
        </Collapse>
      </Box>
    )
  }

  return getGroupView();
}

export default ResultGroup;
