import * as React from 'react';
import ResultGroup from './result_group'

function ResultViewer(props) {
  const headers = props.data.headers ? props.data.headers : ["_key"];
  const groups = props.data.result_groups;

  const getGroupCount = () => {
    if (!props.hasOwnProperty('data')) {
      return 0
    }
    const data = props.data
    if (!data.hasOwnProperty('group_count')) {
      return 0
    }
    return data.group_count
  }

  const getView = () => {
    switch (getGroupCount()) {
      case 0:
        return (
          <div>No Data</div>
        )
      default:
        return groups.map((group) => {
          return (
            <ResultGroup
              headers={headers}
              data={group}
              queryFn={props.queryFn}
            />
          )
        })
      }
  }

  return (
    <div>
      {getView()}
    </div>
  );
}

export default ResultViewer;
