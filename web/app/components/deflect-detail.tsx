import * as React from 'react'






export const DefectDetail = (props: { defect: Defect, attachments: Attachment[], domain: string, project: string,  } ) => {

    // file://[IMAGE_BASE_PATH_PLACEHOLDER]RichContentImage_375600742_1.PNG
    var resultDetail = props.defect.description

    props.attachments.forEach(attachItem => {
        const srcPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
        const destPattern = `/domains/${props.domain}/projects/${props.project}/attachments/${attachItem.id}`;
        resultDetail = resultDetail.replace(srcPattern, destPattern)
    })

    return (
        <>
            <h3>🐞 {props.defect.name} </h3>
            <h4>Detected By: {props["detected-by"]}</h4>
            <div dangerouslySetInnerHTML={{ __html: resultDetail }}></div>

            <h3>🗣 Comments</h3>
            <div style={{}} dangerouslySetInnerHTML={{ __html: props.defect["dev-comments"] }}></div>
        </>
    )
}