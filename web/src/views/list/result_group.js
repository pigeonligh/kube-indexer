import * as React from 'react';
import { Box, Button, Card, Collapse } from '@mui/material';
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
          queryFn={props.queryFn}
        />
      )
    }
    return (
      <Card sx={{
        marginBottom: '5px'
      }}>
        <Card 
          sx={{
            background: '#eeeeee'
          }}
          onClick={() => setCollapseOpen(!collapseOpen)}
        >
          <IconButton
            aria-label="expand row"
            size="large"
          >
            {collapseOpen ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
          <Button
            size="large"
            sx={{textTransform: 'none'}}
          >
            {props.data.name} ({props.data.count})
          </Button>
        </Card>
        
        <Box>
          <Collapse in={collapseOpen} timeout="auto" unmountOnExit sx={{
            padding: '5px',
            background: '#efefef'
          }}>
            <ResultTable 
              name={props.data.name}
              headers={props.headers}
              data={props.data.items}
              queryFn={props.queryFn}
            />
          </Collapse>
        </Box>
      </Card>
    )
  }

  return getGroupView();
}

export default ResultGroup;
