import { RefObject } from 'office-ui-fabric-react';
import * as React from 'react'
import { ContentEditable, DefectComments } from './deflect-comments';
import { DefectEditModal, IDefectEditModal } from './deflect-prop-modal';





export const DefectDetail = (props: DefectPageProps) => {


    // file://[IMAGE_BASE_PATH_PLACEHOLDER]RichContentImage_375600742_1.PNG
    var resultDetail = props.data.defect.description

    props.data.attachment.forEach(attachItem => {
        const srcPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
        const destPattern = `/domains/${props.data.domain}/projects/${props.data.project}/attachments/${attachItem.id}`;
        resultDetail = resultDetail.replace(srcPattern, destPattern)
    })

    const dialogUpdateDefect = () => {
        propertyDialogRef.current.openModal();
    }

    const propertyDialogRef = React.createRef<IDefectEditModal>()
    
    const doSaveDefectProps = (val):Promise<any> => {
        var ret = new Promise<any>((resolve, reject) => {
            resolve('done')
        });
        return ret;
    }
    const onDefectEditModalDismiss = () => {
    }


    const ElmSelf = () => (
        <>
            <DefectEditModal
                ref={propertyDialogRef} 
                data={props}  
                onSave={doSaveDefectProps}
                onDismissed={onDefectEditModalDismiss}   ></DefectEditModal>
            <div className="d-flex justify-content-between">
                <h3>🐞 {props.data.defect.name} </h3>
                <div>
                    <button onClick={dialogUpdateDefect}>✎</button>
                </div>
            </div>
            <div className="panel-footer d-flex justify-content-between">
                <div>
                    <div><b>Detected By:</b> {props.data.defect["detected-by"]} (create on {props.data.defect["creation-time"]} )</div>
                    <div><b>Assigned to:</b> {props.data.defect["owner"]}</div>
                </div>
                <div>
                    <div title="prev - current"><b>Status:</b> {props.data.defect["user-46"]}</div>
                    <div ><b>Last modified:</b> {props.data.defect["last-modified"]}</div>
                </div>
            </div>
            <hr/>
            
            <div dangerouslySetInnerHTML={{ __html: resultDetail }}></div>

            <DefectComments data={props.data}></DefectComments>

            
        </>)

        const elmSelf = ElmSelf();
    

    return (<>
        <div id="selfContainer">
            <ElmSelf/>
        </div>
    </>);
}




