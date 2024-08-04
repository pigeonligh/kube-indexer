import * as React from 'react';
import { Box, Typography } from "@mui/material";
import JsonView from '@uiw/react-json-view';
import { githubLightTheme } from '@uiw/react-json-view/githubLight';

const style = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  bgcolor: 'background.paper',
  boxShadow: 24,
  p: 4,
  width: 800,
};

function RawViewer(props) {
  const [data, setData] = React.useState({});

  React.useEffect(() => {
    if (props.open) {
      fetch('/api/resource/'+props.kind+"/"+props.resourcekey+"?raw=true", {
        method: 'GET'
      }).then(
        (response) => response.json()
      ).then((data) => {
        setData(data)
      }).catch((err) => {
        console.log(err.message);
      });
    }
  }, [props.open, props.kind, props.resourcekey])

  return (
    <Box sx={style}>
      <Typography>
        <div style={{
          overflow: 'scroll',
          maxHeight: '400px',
        }}>
          <JsonView 
            value={data}
            style={githubLightTheme}
            displayDataTypes={false}
            collapsed={3}
            shortenTextAfterLength={0}
            objectSortKeys={true}
          />
        </div>
      </Typography>
    </Box>
  );
}

export default RawViewer;
