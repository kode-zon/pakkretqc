import { Fabric } from '@fluentui/react'
import * as React from 'react'
import { render } from 'react-dom'
import { Attachments, DefectDetail } from './components'

const DefectPage = (props: DefectPageProps) => {

    return (
        <Fabric>
            <div className="defect-container">
                <div className="defect-detail-container" >
                    <DefectDetail data={props.data} />
                </div>
                <div className="defect-attachments-container">
                    <Attachments domain={props.data.domain} project={props.data.project} attachments={props.data.attachment} defectId={props.data.defect.id} />
                </div>
            </div>
        </Fabric>
    )
}


render(<DefectPage data={window.__DATA__}></DefectPage>, document.getElementById("pakkretqc-root"))