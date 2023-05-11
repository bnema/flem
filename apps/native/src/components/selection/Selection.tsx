import React from 'react'
import { StyleSheet, View, Text } from 'react-native'
import { COLORS } from '../../utils/Params'

type SelectionProps = {
    type: string,
}

export default function Selection({ type }: SelectionProps) {
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