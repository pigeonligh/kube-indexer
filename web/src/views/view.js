import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import CssBaseline from '@mui/material/CssBaseline';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import QueryView from './list/query_view';
import Results from './list/results';
import ResultInfo from "./list/result_info";
import ResultSwitch from './list/result_switch';

const drawerWidth = 320;

function MainView(props) {
  const [results, setResults] = React.useState([]);
  const [current, setCurrent] = React.useState(-1);
  const [kinds, setKinds] = React.useState([]);
  const [queryKind, setQueryKind] = React.useState("");

  const [from, setFrom] = React.useState("");
  const [filter, setFilter] = React.useState("");
  const [groupBy, setGroupBy] = React.useState("");

  React.useEffect(() => {
    fetch('/api/kinds', {
      method: 'GET'
    }).then(
      (response) => response.json()
    ).then((data) => {
      setKinds(data)
    }).catch((err) => {
      console.log(err.message);
    });
  }, [props])

  const queryFn = (kind, from, filter, groupBy) => {
    fetch('/api/resource/'+kind, {
      method: 'POST',
      body: JSON.stringify({
        from: from,
        filter: filter,
        group_by: groupBy,
      })
    }).then(
      (response) => response.json()
    ).then((data) => {
      const arr = [...results];
      arr.push(data)
      setResults(arr)
    }).catch((err) => {
      console.log(err.message);
    });
  }

  React.useEffect(() => {
    setCurrent(results.length-1);
  }, [results]);

  const drawer = (
    <div>
      <Toolbar sx={{ background: '#326CE5', color: 'white' }}>
        <Typography variant="h5" noWrap component="div">
          KubeIndexer
        </Typography>
      </Toolbar>
      <Divider />
      <List>
        <ListItem>
          <QueryView
            queryFn={queryFn}
            kinds={kinds}
            queryKind={queryKind}
            setQueryKind={setQueryKind}

            from={from}
            setFrom={setFrom}
            filter={filter}
            setFilter={setFilter}
            groupBy={groupBy}
            setGroupBy={setGroupBy}
          />
        </ListItem>
      </List>
      <Divider />
      <List>
        <ListItem>
          <ResultInfo 
            results={results}
            current={current}
            setCurrent={setCurrent}

            setKind={setQueryKind}
            setFrom={setFrom}
            setFilter={setFilter}
            setGroupBy={setGroupBy}
          />
        </ListItem>
      </List>
    </div>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar sx={{ background: '#326CE5', justifyContent: 'center' }}>
          <ResultSwitch 
            results={results}
            current={current}
            setCurrent={setCurrent}
          />
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
      >
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - ${drawerWidth}px)` } }}
      >
        <Toolbar />
        <Results
          datas={results}
          current={current}
          queryFn={queryFn}
        />
      </Box>
    </Box>
  );
}

export default MainView;