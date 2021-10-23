import { RefObject } from 'office-ui-fabric-react';
import * as React from 'react'
import { render } from 'react-dom'


export const ContentEditable = (props: { 
                                    html: string, 
                                    onChange: (arg) => void,
                                    className?: string
                                },
                                //refObj: RefObject<HTMLDivElement>
                                ) => {

    const [stage, setStage] = React.useState({
        html: props.html
    })

    var lastHtml = '';

    const refObj = React.createRef<HTMLDivElement>();

    function emitChange() {
        if( refObj && refObj.current ) {
            var curHtml = refObj.current.innerHTML;
            if (props.onChange && curHtml !== lastHtml) {
                props.onChange({
                    target: {
                        value: curHtml
                    }
                })
            };
            lastHtml = curHtml;
        }
    }

    function doAppendTimeStamp() {
        console.debug("doAppendTimeStamp");
        if( refObj && refObj.current ) {
            var curHtml = refObj.current.innerHTML;

            curHtml += "<p><span>user name  yyyy-mm-dd HH:MM:ss.SSS<span><p>";

            if (props.onChange) {
                props.onChange({
                    target: {
                        value: curHtml
                    }
                })
            };
            lastHtml = curHtml;
            setStage({ html: lastHtml});
        }
    }

    return (
        <>
            <div className="toolbar-panel">
                <button onClick={doAppendTimeStamp}>append timestamp</button>
            </div>
            <div ref={refObj}
                className={props.className}
                onInput={emitChange}
                onBlur={emitChange}
                contentEditable
                dangerouslySetInnerHTML={{__html: stage.html}}
            >
            </div>
        </>
    )
}


export const DefectDetail = (props: { defect: Defect, attachments: Attachment[], domain: string, project: string,  } ) => {


    // file://[IMAGE_BASE_PATH_PLACEHOLDER]RichContentImage_375600742_1.PNG
    var resultDetail = props.defect.description

    props.attachments.forEach(attachItem => {
        const srcPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
        const destPattern = `/domains/${props.domain}/projects/${props.project}/attachments/${attachItem.id}`;
        resultDetail = resultDetail.replace(srcPattern, destPattern)
    })

    const [stage, setStage] = React.useState({
        commentMode: 'visual',
        comments: props.defect["dev-comments"] + "## sample ##"
    })

    
    var editingComments = stage.comments;

    const commentBtn = (evn) => {
        console.debug("commentBtn invoked", evn);
    }

    const toggleCommentMode = (mode) => {
        console.debug("toggleCommentMode", mode);
        setStage({ 
            commentMode: mode,
            comments: editingComments
        });
    }
    const onContentChange = (arg) => {
        editingComments = arg.target.value;
    }
    const submitComment = () => {
        console.debug("submitComment invoked");
    }


    const ElmSelf = () => (
        <>
            <h3>🐞 {props.defect.name} </h3>
            <h4>Detected By: {props["detected-by"]}</h4>
            <div dangerouslySetInnerHTML={{ __html: resultDetail }}></div>

            <div className="panel-container">
                <div className="panel-header">
                    
                    <ul className="nav nav-tabs navbar-right">
                        <li className="nav-item mr-auto">
                            <h3 onClick={commentBtn}>
                                🗣 Comments 
                            </h3>
                        </li>
                        <li className="nav-item">
                            <a  className={`nav-link ${stage.commentMode==='visual'?'active':null}`}
                                onClick={() => toggleCommentMode('visual')}>
                                visual
                            </a>
                        </li>
                        <li className="nav-item">
                            <a  className={`nav-link ${stage.commentMode==='raw'?'active':null}`}
                                onClick={() => toggleCommentMode('raw')}>
                                raw http
                            </a>
                        </li>
                    </ul>
                </div>
                <div id="defects-comments-area" className="panel-content">
                    { stage.commentMode==='raw'? 
                        <textarea defaultValue={editingComments} onChange={onContentChange} className="comments-area"></textarea> 
                        : 
                        <ContentEditable html={stage.comments} className="comments-area" onChange={onContentChange}></ContentEditable>
                    }
                </div>
                <div className="panel-footer d-flex justify-content-between">
                    <button onClick={() => submitComment()}>submit</button>

                    <div className="d-flex">
                        <button onClick={() => toggleCommentMode('visual')}>visual</button>
                        <button onClick={() => toggleCommentMode('raw')}>raw http</button>
                    </div>
                </div>
            </div>
        </>)

        const elmSelf = ElmSelf();
    

    return (<>
        <div id="selfContainer">
            <ElmSelf/>
        </div>
    </>);
}




