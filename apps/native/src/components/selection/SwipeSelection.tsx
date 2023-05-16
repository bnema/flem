import React from 'react'
import { StyleSheet, View, Text } from 'react-native'
import { COLORS } from '../../utils/Params'

type Colors = {
    BACKGROUND: string;
    UI: string;
    ICON: string;
    BUTTON: string;
    YUP: string;
    NOPE: string;
    VIEW: string;
}

export default function SwipeSelection({ type }: { type: keyof Colors }) {
    const color = COLORS[type]

    return (
        <View style={[styles.container, { borderColor: color }]} >
            <Text style={[styles.text, { color: color }]} >{type}</Text>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        borderWidth: 6,
        paddingHorizontal: 10,
        borderRadius: 5,
        backgroundColor: 'rgba(0,0,0,.2)'
    },
    text: {
        fontSize: 40,
        fontWeight: 'bold',
        textTransform: 'uppercase',
    }
})