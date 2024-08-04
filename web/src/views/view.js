import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import CssBaseline from '@mui/material/CssBaseline';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import MenuIcon from '@mui/icons-material/Menu';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import QueryView from './list/query_view';
import Results from './list/results';
import ResultTab from "./list/result_tab";

const drawerWidth = 240;

function MainView(props) {
  const { window } = props;
  const [mobileOpen, setMobileOpen] = React.useState(false);
  const [isClosing, setIsClosing] = React.useState(false);

  const [results, setResults] = React.useState([]);
  const [current, setCurrent] = React.useState(-1);
  const [kinds, setKinds] = React.useState([]);
  const [queryKind, setQueryKind] = React.useState("");

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

  const queryFn = (kind, filter, groupBy) => {
    fetch('/api/resource/'+kind, {
      method: 'POST',
      body: JSON.stringify({
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

  const handleDrawerClose = () => {
    setIsClosing(true);
    setMobileOpen(false);
  };

  const handleDrawerTransitionEnd = () => {
    setIsClosing(false);
  };

  const handleDrawerToggle = () => {
    if (!isClosing) {
      setMobileOpen(!mobileOpen);
    }
  };

  const drawer = (
    <div>
      <Toolbar />
      <Divider />
      <List>
        <ListItem>
          <QueryView
            queryFn={queryFn}
            kinds={kinds}
            queryKind={queryKind}
            setQueryKind={setQueryKind}
          />
        </ListItem>
      </List>
      <Divider />
      <List>
        <ListItem>
          <ResultTab 
            results={results}
            current={current}
            setCurrent={setCurrent}
          />
        </ListItem>
      </List>
    </div>
  );

  // Remove this const when copying and pasting into your project.
  const container = window !== undefined ? () => window().document.body : undefined;

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
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            KubeIndexer Query
          </Typography>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
        aria-label="mailbox folders"
      >
        {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
        <Drawer
          container={container}
          variant="temporary"
          open={mobileOpen}
          onTransitionEnd={handleDrawerTransitionEnd}
          onClose={handleDrawerClose}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
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