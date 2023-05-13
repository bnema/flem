import React, { useEffect, useState } from 'react'
import { StyleSheet, View, Text } from 'react-native'
import GenreSelection from '../selection/GenreSelection'
import genresDB from '../../utils/genresDB'
import { COLORS } from '../../utils/Params'

export default function Genres() {

    return (
        <View>
            <Text style={styles.text}>What should you watch?</Text>
            <View style={styles.genresContainer}>
                {genresDB.map(({ id, name }) => {
                    return <GenreSelection key={id} genre={name} />
                })}
            </View>
        </View>
    )
}

const styles = StyleSheet.create({
    text: {
        left: 20,
        fontSize: 20,
        fontWeight: '500',
        color: COLORS.VIEW,
        marginBottom: 25,
    },
    genresContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        columnGap: 15,
        rowGap: 20,
        paddingHorizontal: 20,
    }
})