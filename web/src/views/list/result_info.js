import * as React from 'react';
import { CodeBlock, dracula } from "react-code-blocks";

function ResultInfo(props) {
  const results = props.results;
  const current = props.current;

  const setKind = props.setKind;
  const setFrom = props.setFrom;
  const setFilter = props.setFilter;
  const setGroupBy = props.setGroupBy;

  React.useEffect(() => {
    if (current < 0 || current >= results.length) {
      return
    }
    const list_param = results[current].list_param;
    setKind(results[current].kind)
    setFrom(list_param.from)
    setFilter(list_param.filter)
    setGroupBy(list_param.group_by)
  }, [results, current, setKind, setFrom, setFilter, setGroupBy])

  const viewFrom = (list_param) => {
    if (list_param.from === "") {
      return
    }
    return (
      <div>
        <span>From: </span>
        <CodeBlock
          language="go"
          text={list_param.from}
          codeBlock
          theme={dracula}
          showLineNumbers={false}
        />
      </div>
    )
  }

  const viewFilter = (list_param) => {
    if (list_param.filter === "") {
      return
    }
    return (
      <div>
        <span>Filter: </span>
        <CodeBlock
          language="go"
          text={list_param.filter}
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
        <span>Group By: </span>
        <CodeBlock
          language="go"
          text={list_param.group_by}
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
        <span>Kind: {results[current].kind}</span>
        {viewFrom(results[current].list_param)}
        {viewFilter(results[current].list_param)}
        {viewGroupBy(results[current].list_param)}
      </div>
    )
  }

  return (
    <div style={{
        maxWidth: '100%'
    }}>
      {currentInfo()}
    </div>
  );
}

export default ResultInfo;
