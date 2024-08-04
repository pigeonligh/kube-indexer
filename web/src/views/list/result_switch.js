import IconButton from '@mui/material/IconButton';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';

function ResultSwitch(props) {
  const results = props.results;
  const current = props.current;
  const setCurrent = props.setCurrent;

  return (
    <div style={{
        maxWidth: '100%'
    }}>
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
    </div>
  );
}

export default ResultSwitch;
