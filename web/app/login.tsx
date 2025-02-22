
import { Dropdown, Fabric, PrimaryButton, Text, TextField } from '@fluentui/react';
import * as React from 'react';
import { render } from 'react-dom';
import { MadeWithLove } from './components';



interface LoginPageData {
    username: string
    domains: Array<{ name: string }>
    redirectUrl?: string
    errorMessage?: string
}



export const LoginPage = (props: { data: LoginPageData }) => {
    const [domain, setDomain] = React.useState<string>("")
    return (
        <Fabric className="content">
            <div className="login-form-container" style={{}}>
                <h1>PakkretQC</h1>
                <h2>Before you begin please sign-in 🚪</h2>
                <form className="login-form" method="POST" action="/login">
                    <TextField disabled={!!props.data.username} defaultValue={props.data.username} name="username" placeholder="username" className="text-field" id="username"></TextField>
                    <TextField disabled={!!props.data.username} name="password" placeholder="password" type="password" className="text-field" id="password"></TextField>
                    <input type="hidden" name="redirect" value={props.data.redirectUrl} />
                    {
                        props.data.errorMessage ?
                            <Text variant="small" style={{color: "#a80000" }}>{props.data.errorMessage}</Text> : null
                    }
                    <div style={{ marginTop: 8, textAlign: "right" }}>
                        <PrimaryButton disabled={!!props.data.username} as="input" type="submit">Authenticate</PrimaryButton>
                    </div>
                    {
                        props.data.username ? (
                            <>
                                <Dropdown
                                    label="Go on pick your Domain and then please come-in 🔓"
                                    selectedKey={domain}
                                    placeHolder={"select one here"}
                                    onChange={(event, item) => setDomain(item.data)}
                                    options={props.data.domains.map(domain => {
                                        return {
                                            key: domain.name,
                                            text: domain.name,
                                            data: domain.name,
                                        }
                                    })}
                                ></Dropdown>
                                <input type="hidden" name="currentDomain" value={domain} />
                                <div className="domain-confirm-actions" style={{ textAlign: 'right' }}>
                                    <PrimaryButton as="input" name="action" value="cancel" type="submit">Cancel</PrimaryButton>
                                    <PrimaryButton disabled={!domain} as="input" name="action" value="proceed" type="submit">Proceed</PrimaryButton>
                                </div>
                            </>
                        ) : null
                    }
                </form>
                <MadeWithLove></MadeWithLove>
            </div>

        </Fabric>
    )
}
render(<LoginPage data={window.__DATA__}></LoginPage>, document.getElementById("pakkretqc-root"))