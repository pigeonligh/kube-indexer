import { CopyBlock, dracula } from "react-code-blocks";
import IconButton from '@mui/material/IconButton';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';

function ResultTab(props) {
  const results = props.results;
  const current = props.current;
  const setCurrent = props.setCurrent;

  const viewFilter = (list_param) => {
    if (list_param.filter === "") {
      return
    }
    return (
      <div>
        <p>Filter: </p>
        <CopyBlock
          language="go"
          text={results[current].list_param.filter}
          codeBlock
          theme={dracula}
          showLineNumbers={false}
        />
      </div>
    )
  }

  const viewGroupBy = (list_param) => {
    if (list_param.group_by === "") {
      return
    }
    return (
      <div>
        <p>Group By: </p>
        <CopyBlock
          language="go"
          text={results[current].list_param.group_by}
          codeBlock
          theme={dracula}
          showLineNumbers={false}
        />
      </div>
    )
  }

  const currentInfo = () => {
    if (current < 0 || current >= results.length) {
      return
    }
    return (
      <div>
        <p>Kind: {results[current].kind}</p>
        {viewFilter(results[current].list_param)}
        {viewGroupBy(results[current].list_param)}
      </div>
    )
  }

  return (
    <div style={{
        maxWidth: '100%'
    }}>
      <div>
      </div>
      <p>History: 
        <IconButton
          aria-label="expand row"
          size="small"
          onClick={() => {
            if (current > 0) {
              setCurrent(current-1);
            }
          }}
        >
          <KeyboardArrowLeftIcon />
        </IconButton>
        {current + 1} / {results.length}
        <IconButton
          aria-label="expand row"
          size="small"
          onClick={() => {
            if (current < results.length - 1) {
              setCurrent(current+1);
            }
          }}
        >
          <KeyboardArrowRightIcon />
        </IconButton>
      </p>
      {currentInfo()}
    </div>
  );
}

export default ResultTab;
