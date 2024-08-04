import ResultViewer from "./result_viewer";

function Results(props) {
  const results = props.datas;
  const current = props.current;

  const getCurrentViewer = () => {
    if (current < 0 || current >= results.length) {
      return
    }
    return (
      <ResultViewer data={results[current]} queryFn={props.queryFn}/>
    )
  }

  return (
    <div>
      {getCurrentViewer()}
    </div>
  );
}

export default Results;
