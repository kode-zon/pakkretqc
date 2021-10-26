import { RefObject, Modal, IDragOptions, ContextualMenu, mergeStyleSets, getTheme, FontWeights, IconButton, IIconProps, ComboBox, IComboBox, IComboBoxOption } from 'office-ui-fabric-react';
import * as React from 'react'

export interface IDefectEditModal {
    openModal: () => void;
}

export const DefectEditModal = React.forwardRef((props: { 
                                    data:DefectPageProps, 
                                    onSave: (data) => Promise<any>, 
                                    onDismissed: () => void 
                                },
                                ref: RefObject<IDefectEditModal>) => {

    const dragOptions: IDragOptions = {
        moveMenuItemText: 'Move',
        closeMenuItemText: 'Close',
        menu: ContextualMenu,
    }
    const theme = getTheme()
    const contentStyles = mergeStyleSets({
        container: {
            "min-width": "calc(50%) !important"
        },
        header: [
            theme.fonts.xLarge,
            {
                flex: '1 1 auto',
                //borderTop: `4px solid ${theme.palette.themePrimary}`,
                backgroundColor: theme.palette.themePrimary,
                color: theme.palette.neutralPrimary,
                display: 'flex',
                alignItems: 'center',
                fontWeight: FontWeights.semibold,
                padding: '12px 12px 14px 24px'
            }
        ],
        body: {
            flex: '4 4 auto',
            padding: `0 24px 24px 24px`,
            overflowY: 'auto',
            selectors: {
                p: { margin: '14px 0' },
                'p:first-child': { marginTop: 0 },
                'p:last-child': { marginBottom: 0 }
            }
        },
        footer: {
            position: 'absolute',
            bottom: '0',
            width: '100%',
            'align-items': 'center',
            'justify-content': 'center'

        }
    })
    const iconButtonStyles = {
        root: {
            color: theme.palette.neutralPrimary,
            marginLeft: 'auto',
            marginTop: '4px',
            marginRight: '2px'
        }
    }
    const cancelIcon:IIconProps = { iconName: 'Cancel' }

    const [isModalOpen, setIsModalOpen] = React.useState(false)
    const [modalState, setModalState] = React.useState({
        almUsedList: {
            statusOptions: [],
            usernameOptions: []
        },
        defectProp: { ...props.data.data.defect },  //clone value, not use =
        structForSave: { 
            Fields: {}
        }
    })
    var editingModalState = {
        ...modalState
    }

    React.useImperativeHandle(ref, () => ({
        openModal() {

            //reset value back to props.data.data.defect.
            editingModalState = { 
                ...modalState,
                defectProp: { ...props.data.data.defect },
                structForSave: {
                    Fields: {}
                }
            };
            setModalState(editingModalState);
            setIsModalOpen(true);

            if(modalState.almUsedList.statusOptions.length < 1) {
                loadListOfValue();
            }
            if(modalState.almUsedList.usernameOptions.length < 1) {
                loadListOfUser();
            }
        }
    }))

    
    const saveBtnRef = React.createRef<HTMLButtonElement>();
    const setDisabledSaveBtn = (value:boolean) => {
        if(saveBtnRef.current != null) {
            saveBtnRef.current.disabled = value;
        }
    }

    const loadListOfValue = () => {

        var pageData = props.data.data;
        const targetUrl = `/domains/${pageData.domain}/projects/${pageData.project}/customization/used-lists`;
        const reqMetadata = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        }

        setDisabledSaveBtn(true);
        fetch(targetUrl, reqMetadata)
            .then(res => res.json())
            .then((resp:{lists:ALMUsedListEntry[]} ) => {
                console.debug("loadListOfValue success", resp);
                var statusList = resp.lists.find(item => { return (item.Id == 4) }).Items;
                editingModalState = { ...modalState };
                editingModalState.almUsedList.statusOptions = statusList.map(item => {
                    return { key: item.value, text: item.value }
                })
                setModalState(editingModalState);
                setDisabledSaveBtn(false);
            })
            .catch(reason => {
                setDisabledSaveBtn(false);
                console.error("error at loadListOfValue", reason)
                alert("error:"+reason)
            });
        console.debug("loadListOfValue invoked");
    }

    
    const loadListOfUser = () => {

        var pageData = props.data.data;
        const targetUrl = `/domains/${pageData.domain}/projects/${pageData.project}/customization/users`;
        const reqMetadata = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        }
        setDisabledSaveBtn(true);
        fetch(targetUrl, reqMetadata)
            .then(res => res.json())
            .then( (resp:{users:[]}) => {
                console.debug("loadListOfUser success", resp)
                editingModalState = { ...modalState };
                editingModalState.almUsedList.usernameOptions = resp.users.map(item => {
                    return { key: item["Name"], text: item["FullName"] }
                })
                setModalState(editingModalState);
                setDisabledSaveBtn(false);
            })
            .catch(reason => {
                setDisabledSaveBtn(false);
                console.error("error at loadListOfUser", reason)
                alert("error:"+reason)
            });
        console.debug("loadListOfUser invoked");
    }


    
    const doSave = () => {

        var pageData = props.data.data;
        const putUrl = `/domains/${pageData.domain}/projects/${pageData.project}/defects/${pageData.defect.id}`;
        const postBody = editingModalState.structForSave;
        // {
        //     Fields: {
        //         "dev-comments": editingComments
        //     }
        // }
        const reqMetadata = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postBody)
        }

        
        console.debug("doSave", postBody);
        setDisabledSaveBtn(true);

        fetch(putUrl, reqMetadata)
            .then(res => {
                if(res.ok) {
                    let resp = res.json();

                    console.log("save response", resp)
                    alert("done");
                    setDisabledSaveBtn(false);
                }
                res.text().then(text => {
                    setDisabledSaveBtn(false);
                    alert("error:"+text)
                });
            })
            .catch(reason => {
                setDisabledSaveBtn(false);
                alert("error:"+reason)
            });
    }
    const doCancel = () => {
        // editingModalState.isModalOpen = false;
        // setModalState(editingModalState);
        setIsModalOpen(false);
    }

    const statusErrorMsg = ():string => {
        if(editingModalState.almUsedList.statusOptions.length < 1) {
            return "list of status doesn't load yet !!";
        }
        return undefined;
    }
    const assignErrorMsg = ():string => {
        if(modalState.almUsedList.usernameOptions.length < 1) {
            return "list of user doesn't load yet !!";
        }
        return undefined;
    }

    const statusChanged = (event:React.FormEvent<IComboBox>, option?: IComboBoxOption, index?: number, value?: string) => {
        editingModalState.defectProp.status = option.key as string;
        editingModalState.structForSave.Fields["status"] = option.key as string;
        setModalState(editingModalState);
        console.debug("statusChanged", editingModalState.defectProp);
    }
    const assignChanged = (event:React.FormEvent<IComboBox>, option?: IComboBoxOption, index?: number, value?: string) => {
        editingModalState.defectProp.owner = option.key as string;
        editingModalState.structForSave.Fields["owner"] = option.key as string;
        setModalState(editingModalState);
        console.debug("assignChanged", editingModalState.defectProp);
    }

    return (
        <>
            <Modal
                titleAriaId="Edit Properties"
                dragOptions={dragOptions}
                isBlocking={true}
                isOpen={isModalOpen}
                onDismissed={props.onDismissed}
                containerClassName={contentStyles.container}>

                <div className={contentStyles.header}>
                    <span>title</span>
                    <IconButton
                        styles={iconButtonStyles}
                        iconProps={cancelIcon}
                        ariaLabel="Close popup modal"
                        onClick={doCancel}>
                    </IconButton>
                </div>

                <div className={contentStyles.body}>
                    <div className="row">
                        <div className="col-6">
                            <ComboBox label="status"
                                //styles={ {callout: { 'max-height': 'calc(60vh)'} } }
                                options={modalState.almUsedList.statusOptions}
                                selectedKey={editingModalState.defectProp.status}
                                errorMessage={statusErrorMsg()}
                                onChange={statusChanged}
                                ></ComboBox>
                        </div>
                        <div className="col-6">
                            <ComboBox label="assign to"
                                options={modalState.almUsedList.usernameOptions}
                                selectedKey={editingModalState.defectProp.owner}
                                errorMessage={assignErrorMsg()}
                                onChange={assignChanged}
                                ></ComboBox>
                        </div>
                    </div>
                </div>

                <div className={contentStyles.footer + " d-flex"}>
                    <div>
                        <button onClick={doSave} ref={saveBtnRef}>save</button>
                        <button onClick={doCancel}>cancel</button>
                    </div>
                </div>
            </Modal>
        </>
    )
})