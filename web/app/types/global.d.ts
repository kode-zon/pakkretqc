
declare interface Window { __DATA__: any }

declare type ContentWrapperMode = 'wrap'|'unwrap'

declare interface ALMUsedListItemEntry {
    "LogicalName": string,
    "value": string
}
declare interface ALMUsedListEntry {
    "Name": string,
    "Id": number,
    "LogicalName": string,
    "Items": ALMUsedListItemEntry[]
}

declare interface DefectPageProps {
    data: {
        defect: Defect
        attachment: Attachment[]
        project: string
        domain: string
        username: string
        userfullname: string
    }
}

declare interface Defect {
    id: number,
    status: string
    "user-46": string // this is extra status (show as [old status] - [current status])
    owner: string
    name: string
    severity: string
    description: string
    "dev-comments": string
    "last-modified": string
    "creation-time": string
    "detected-by": string

    url: string
}

declare interface Attachment {

    type: string;
    "last-modified": string;
    "vc-cur-ver"?: any;
    "vc-user-name"?: any;
    name: string;
    "file-size": number;
    "ref-subtype": number;
    description?: any;
    id: number;
    "ref-type": string;
    entity: {
        id: string
        type: string
    };
}