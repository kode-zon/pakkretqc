import { RefObject } from 'office-ui-fabric-react';
import * as React from 'react'

export const ContentEditable = React.forwardRef((props: { 
                                    html: string, 
                                    onChange: (arg) => void,
                                    className?: string,
                                    username: string,
                                    userfullname: string,
                                    contentWrapper?: (mode: ContentWrapperMode, content:string) => string,
                                    mode: string
                                },
                                ref: RefObject<any>
                                ) => {
                                    
    React.useImperativeHandle(ref, () => ({
        appendStamp() {
            doAppendTimeStamp()
        }
    }))

    const [stage, setStage] = React.useState({
        html: props.html
    })

    var lastHtml = '';

    const refObj = React.createRef<HTMLDivElement>();

    function wrapContent(content:string):string {
        return (props.contentWrapper)?props.contentWrapper('wrap',content):content;
    }
    function unwrapContent(content:string):string {
        return (props.contentWrapper)?props.contentWrapper('unwrap',content):content;
    }

    function emitChange() {
        if( refObj && refObj.current ) {
            var curHtml = unwrapContent(refObj.current.innerHTML);
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

            const timestamp = new Date();
            let ye = new Intl.DateTimeFormat('en', { year: 'numeric'}).format(timestamp);
            let mo = new Intl.DateTimeFormat('en', { month: '2-digit'}).format(timestamp);
            let da = new Intl.DateTimeFormat('en', { day: '2-digit'}).format(timestamp);
            let timeStr = new Intl.DateTimeFormat('en', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false}).format(timestamp);
            
            var timestampStr = `${da}/${mo}/${ye} ${timeStr}`   // example "22/10/2021 12:48:25"
            curHtml += `<p><font face="Arial" color="#000055"><span dir="ltr" style="font-size:8pt">________________________________________________<br /></span><span dir="ltr" style="font-size:8pt"><b>${props.userfullname} &lt; ${props.username} &gt; ${timestampStr} </b></span></font><br />-</p>`;
            

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
            {
                (props.mode==='visual')?
                    <div className="toolbar-panel">
                        <button onClick={doAppendTimeStamp}>append timestamp</button>
                    </div>
                :
                    undefined
            }
            
            <div ref={refObj}
                className={props.className}
                onInput={emitChange}
                onBlur={emitChange}
                contentEditable={(props.mode==='visual')?true:false}
                dangerouslySetInnerHTML={{__html: wrapContent(stage.html) }}
            >
            </div>
        </>
    )
})


export interface IDefectComments {
    setCommentMode: (mode) => void;
}

export const DefectComments = (props:  DefectPageProps) => {
    
    const [stage, setStage] = React.useState({
        commentMode: 'view',
        comments: props.data.defect["dev-comments"]
    })
    const contentEditableRef = React.createRef<any>()
    const submitBtnRef = React.createRef<HTMLButtonElement>()

    var editingComments = stage.comments;

    const commentBtn = (evn) => {
        console.debug("commentBtn invoked", evn);

        contentEditableRef.current.appendStamp();
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

        const putUrl = window.location.href;
        // const postBody = {
        //     Fields: [{
        //         Name: "dev-comments",
        //         values: [{
        //             "value": editingComments
        //         }]
        //     }]
        // }
        const postBody = {
            Fields: {
                "dev-comments": editingComments
            }
        }
        const reqMetadata = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postBody)
        }

        submitBtnRef.current.disabled = true

        fetch(putUrl, reqMetadata)
            .then(res => res.json())
            .then(resp => {
                console.log("save response", resp)
                alert("done");
                submitBtnRef.current.disabled = false
            })
            .catch(reason => {
                submitBtnRef.current.disabled = false
                alert("error:"+reason)
            })
            
    }

    const contentWrapperFn = (mode: ContentWrapperMode, content:string):string => {
        var resultContent: string = content;
        if(mode=='wrap') {
            props.data.attachment.forEach(attachItem => {
                const srcPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
                const destPattern = `/domains/${props.data.domain}/projects/${props.data.project}/attachments/${attachItem.id}`;
                resultContent = resultContent.replace(srcPattern, destPattern)
            })
        } else {
            props.data.attachment.forEach(attachItem => {
                const srcPattern = `/domains/${props.data.domain}/projects/${props.data.project}/attachments/${attachItem.id}`;
                const destPattern = "file://[IMAGE_BASE_PATH_PLACEHOLDER]"+attachItem.name;
                resultContent = resultContent.replace(srcPattern, destPattern)
            })
        }

        return resultContent;
    }


    return (
        <>
            <div className="panel-container">
                <div className="panel-header">
                    
                    <ul className="nav nav-tabs navbar-right">
                        <li className="nav-item mr-auto">
                            <h3>
                                🗣 
                                {
                                    (stage.commentMode==='raw' || stage.commentMode==='visual') ?
                                        <button onClick={commentBtn}>Comments</button>
                                    :
                                        <span>Comments</span>
                                }
                            </h3>
                        </li>
                        <li className="nav-item">
                            <a  className={`nav-link ${stage.commentMode==='view'?'active':null}`}
                                onClick={() => toggleCommentMode('view')}>
                                view
                            </a>
                        </li>
                        <li className="nav-item">
                            <a  className={`nav-link ${stage.commentMode==='visual'?'active':null}`}
                                onClick={() => toggleCommentMode('visual')}>
                                ✎ visual
                            </a>
                        </li>
                        <li className="nav-item">
                            <a  className={`nav-link ${stage.commentMode==='raw'?'active':null}`}
                                onClick={() => toggleCommentMode('raw')}>
                                 ✎ http raw
                            </a>
                        </li>
                    </ul>
                </div>
                <div id="defects-comments-area" className="panel-content">
                    { stage.commentMode==='raw'? 
                        <textarea defaultValue={editingComments} onChange={onContentChange} className="comments-area"></textarea> 
                        : 
                        <ContentEditable ref={contentEditableRef} 
                            html={stage.comments} className="comments-area" 
                            onChange={onContentChange} contentWrapper={contentWrapperFn}
                            username={props.data.username}
                            userfullname={props.data.userfullname}
                            mode={stage.commentMode}/>
                    }
                </div>
                <div className="panel-footer d-flex justify-content-between">
                    <button ref={submitBtnRef} onClick={() => submitComment()}>save comments</button>
                </div>
            </div>
        </>
    )
}