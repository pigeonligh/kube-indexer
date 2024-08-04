import JsonView from '@uiw/react-json-view';
import { githubLightTheme } from '@uiw/react-json-view/githubLight';

function ResultObjectViewer(props) {
  const data = JSON.parse(JSON.stringify(props.data));

  delete data._raw;
  delete data._key;
  delete data._kind;
  delete data._resource_version;

  return (
    <div>
      <JsonView 
        value={data}
        style={githubLightTheme}
        displayDataTypes={false}
        collapsed={false}
        shortenTextAfterLength={0}
        objectSortKeys={true}
      />
    </div>
  );
}

export default ResultObjectViewer;
