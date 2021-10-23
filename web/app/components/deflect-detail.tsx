import * as React from 'react'
import { render } from 'react-dom'



export const DefectDetail = (props: { defect: Defect, attachments: Attachment[], domain: string, project: string,  } ) => {

    // file://[IMAGE_BASE_PATH_PLACEHOLDER]RichContentImage_375600742_1.PNG
    var resultDetail = props.defect.description

    props.attachments.forEach(attachItem => {
        const srcPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
        const destPattern = `/domains/${props.domain}/projects/${props.project}/attachments/${attachItem.id}`;
        resultDetail = resultDetail.replace(srcPattern, destPattern)
    })

    
    var comments = props.defect["dev-comments"];
    var commentMode = 'visual'

    console.debug("comments = ", comments);

    const fnUpdate = () => {
        
        render(<ElmComments></ElmComments>, document.getElementById("defects-comments-area"))
    }

    const commentBtn = (evn) => {
        console.debug("commentBtn invoked", evn);
    }

    const toggleCommentMode = (mode) => {
        console.debug("toggleCommentMode", mode);
        commentMode = mode;
        console.debug("commentMode=raw? : ", commentMode==='raw');
        fnUpdate();
    }
    const onContentChange = (arg) => {
        comments = arg.target.value;
        fnUpdate();
    }
    const submitComment = () => {
        console.debug("submitComment invoked");
    }


    const ElmCommentsVisual = () => (
        <div dangerouslySetInnerHTML={{ __html: comments }}></div>
    )
    const ElmCommentsRaw = () => (
        <textarea defaultValue={comments} onChange={onContentChange} className="w-100"></textarea>
    )

    
    const ElmComments = () => (
        <>
            <ElmCommentsVisual/>
            { commentMode==='raw'? <ElmCommentsRaw/> : null }
            <br/>
            <button onClick={() => submitComment()}>submit</button>
        </>      
    )
    

    return (
        <>
            <h3>🐞 {props.defect.name} </h3>
            <h4>Detected By: {props["detected-by"]}</h4>
            <div dangerouslySetInnerHTML={{ __html: resultDetail }}></div>

            <div className="defects-comments-container">
                <h3 onClick={commentBtn}>🗣 Comments 
                    <button onClick={() => toggleCommentMode('visual')}>visual</button>
                    <button onClick={() => toggleCommentMode('raw')}>raw http</button>
                </h3>
                <div id="defects-comments-area">
                    <ElmComments></ElmComments>
                </div>
            </div>
            
        </>
    )
}




