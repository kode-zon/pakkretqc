import { ActionButton, ComboBox, DatePicker, DayOfWeek, IComboBox, IComboBoxOption, IDatePickerStrings, IIconProps, initializeIcons, ISelectableOption, RefObject } from 'office-ui-fabric-react';
import * as React from 'react'
import { ContentEditable, DefectComments } from './deflect-comments';
import { DefectEditModal, IDefectEditModal } from './deflect-prop-modal';



initializeIcons();

const DayPickerString: IDatePickerStrings = {
    months: [
        'January',
        'Febuary',
        'March',
        'April',
        'May',
        'June',
        'July',
        'August',
        'September',
        'Octtober',
        'November',
        'December'
    ],
    shortMonths: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
    days: [],
    shortDays: [],
    goToToday: 'go to today'
}
const TimeOptionsConst = {
    HH: [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
         10,11,12,13,14,15,16,17,18,19,
         20,21,22,23,24],
    MM: [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
         10,11,12,13,14,15,16,17,18,19,
         20,21,22,23,24,25,26,27,28,29,
         30,31,32,33,34,35,36,37,38,39,
         40,41,42,43,44,45,46,47,48,49,
         50,51,52,53,54,55,56,57,58,59],
}

const TimeOptions:{HH:ISelectableOption[], MM:ISelectableOption[]} = {
    HH: TimeOptionsConst.HH.map(item => { return { key:item, text:""+item }}),
    MM: TimeOptionsConst.MM.map(item => { return { key:item, text:""+item }})
}

export interface IDateTimePicker {
    setValue: (date:Date) => void
}

export const DateTimePicker = (props: { 
                                    dateTimeValue:Date|string, 
                                    label?:string,
                                    placeholder?:string,
                                    onChanged: (data:Date) => any
                                }) => {

    const [stage, setStage] = React.useState({
        dateValue: (props.dateTimeValue!=null)?new Date(props.dateTimeValue):null
    });

    // if(ref != null) {
    //     React.useImperativeHandle(ref, () => ({
    //         setValue(date:Date) {
    //             doSetDateValue(date);
    //         }
    //     }));
    // }

    const doSetDateValue = (modDate:Date) => {
        // console.debug("doSetDateValue", modDate);
        // console.debug("doSetDateValue", stage);
        setStage({ dateValue:modDate });
        // console.debug("doSetDateValue", stage);
        props.onChanged(modDate);
    }

    const dateValueChanged = (value:Date) => {
        let modDate = value;
        if(stage.dateValue!=null) {
            modDate = stage.dateValue;
            modDate.setDate(value.getDate());
            modDate.setMonth(value.getMonth());
            modDate.setFullYear(value.getFullYear());
        }
        
        doSetDateValue(modDate);
    }
    const selectHHChanged = (event:React.FormEvent<IComboBox>, option?: IComboBoxOption, index?: number, value?: string) => {
        let modDate = stage.dateValue;
        modDate.setHours(option.key as number);
        doSetDateValue(modDate);
    }
    
    const selectMMChanged = (event:React.FormEvent<IComboBox>, option?: IComboBoxOption, index?: number, value?: string) => {
        let modDate = stage.dateValue;
        option.key as number;
        modDate.setMinutes(option.key as number);
        doSetDateValue(modDate);
    }

    const deleteDate = () => {
        doSetDateValue(null);
    }

    const trashIcon:IIconProps = {
        iconName: 'trash'
    }


    return (<>
        <div className="row">
            <div className="col-12">
                <DatePicker label={props.label}
                    value={stage.dateValue!}
                    firstDayOfWeek={DayOfWeek.Sunday}
                    strings={DayPickerString}
                    placeholder={props.placeholder}
                    onSelectDate={dateValueChanged}
                    ></DatePicker>
            </div>
            
            <div className="col-12">
                <div className="row">
                    <div className="col-11">
                        <div className="row">
                            <div className="col-6 mr-0 pr-1">
                                <ComboBox 
                                    disabled={stage.dateValue==null}
                                    options={TimeOptions.HH}
                                    selectedKey={stage.dateValue?.getHours()}
                                    onChange={selectHHChanged}
                                ></ComboBox>
                            </div>
                            <div className="col-6 ml-0 pl-1">
                                <ComboBox 
                                    disabled={stage.dateValue==null}
                                    options={TimeOptions.MM}
                                    selectedKey={stage.dateValue?.getMinutes()}
                                    onChange={selectMMChanged}
                                ></ComboBox>
                            </div>
                        </div>
                    </div>
                    <div className="col-1 mr-auto">
                        <ActionButton iconProps={trashIcon} onClick={deleteDate}></ActionButton>
                    </div>
                </div>
            </div>
            
            <div className="col-12">debug :: {""+stage.dateValue!}</div>
        </div>
    </>);
}




