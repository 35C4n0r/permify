import React from "react";
import AppRouter from "@router/router";
import {Layout} from "@utility/context/Layout"
import {ConfigProvider, theme} from 'antd';
import { Analytics } from '@vercel/analytics/react';
function App() {
    return (
        <ConfigProvider
            theme={{
                token: {
                    colorPrimary: '#6318FF',
                    fontSizeBase: '14px',
                    borderRadius: '2px',
                },
                components: {
                    divider: {
                        backgroundColor: '#DEDBE6',
                    }
                },
                algorithm: theme.darkAlgorithm,
            }}
        >
            <Layout>
                <AppRouter/>
                <Analytics />
            </Layout>
        </ConfigProvider>
    );
}

export default App
