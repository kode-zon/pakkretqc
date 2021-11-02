import { DetailsList, IColumn, Link, SelectionMode } from '@fluentui/react'
import * as React from 'react'

const columns: (domain: string, project: string) => IColumn[] = (domain, project) => [
    {
        key: 'col1',
        name: "ID",
        fieldName: "id",
        minWidth: 50,
        maxWidth: 50,
    },
    {
        key: 'col2',
        name: "Name",
        fieldName: "name",
        minWidth: 255,
        maxWidth: 455,
        data: 'string',
        isRowHeader: true,
        onRender: (item, idx, col) => {
            return <Link target="popup" download={item['name']} href={`/domains/${domain}/projects/${project}/attachments/${item.id}`}>{item['name']}</Link>
        }
    },
    {
        key: 'col22',
        name: 'Size',
        fieldName: 'file-size',
        minWidth: 120,
        isResizable: true,
    },

]

export const Attachments = (props: { domain: string, project: string, defectId: number, attachments: Attachment[] }) => {

    const [state, setState] = React.useState({theFile:null, desc:"desc of file"})

    const changeFile = (e) => {
        setState({theFile: e.target.files[0], desc:state.desc})
    }
    const changeDesc = (e) => {
        setState({theFile: state.theFile, desc:e.target.value})
    }
    const uploadFile = () => {
        let formData = new FormData();
        formData.append("description", `${state.desc}`);
        formData.append("filename", state.theFile.name);
        formData.append("file", state.theFile, state.theFile.name);
        
        let targetUrl =`/domains/${props.domain}/projects/${props.project}/defects/${props.defectId}/attachments`
        fetch(
            targetUrl,
            {
                method: 'POST',
                body: formData
            }
        ).then((resp) => {
            console.debug("uploadFile response")
            if(resp.ok) {
                resp.json().then(result => {
                    console.log("Attach success:", result);
                });
                return;
            }
            resp.text().then(text => {
                alert("error:"+text)
            });
        })
        .catch(reason => {
            alert("error:"+reason)
        });
    }

    return (
        <>
            <h3>📑 Attachments</h3>
            <div >
                <DetailsList
                    selectionMode={SelectionMode.none}
                    columns={columns(props.domain, props.project)}
                    items={props.attachments}
                />
                <hr/>
                <div>
                    <input type="file" onChange={changeFile} className="border"></input>
                    <br/>
                    note : <input type="text" onChange={changeDesc} className="border"></input>
                    <br/>
                    <button onClick={uploadFile} disabled={state.theFile==null}>Upload</button>
                </div>
            </div>
        </>
    )
}