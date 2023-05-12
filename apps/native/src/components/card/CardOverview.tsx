import React from 'react'
import { StyleSheet, View, Text } from 'react-native'
import { CARD } from '../../utils/Params'

type MoviesProps = {
    title: string,
    genre: string[],
    overview: string,
    date: string,
}

export default function CardOverview({ title, genre, overview, date, }: MoviesProps) {
    return (
        <View style={styles.container} >
            <Text style={styles.title}>{title}</Text>
            <Text style={styles.genre}>{genre.join(' - ')}</Text>
            <Text style={styles.overview}>{overview.replace(/^((?:\S+\s+){0,40}).*/, "$1") + '...'}</Text>
            <Text style={styles.date}>Release date: {date.slice(0, 4)}</Text>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        width: CARD.WIDTH,
        height: 280,
        borderBottomLeftRadius: CARD.BORDER_RADIUS,
        borderBottomRightRadius: CARD.BORDER_RADIUS,
        backgroundColor: '#fff',
        borderWidth: 6,
        borderColor: '#000',
        paddingHorizontal: 20,
        paddingVertical: 40,
    },
    title: {
        fontSize: 20,
        fontWeight: 'bold'
    },
    genre: {
        fontWeight: '300',
        color: '#808080',
    },
    overview: {
        fontStyle: 'italic',
        marginVertical: 2,
    },
    date: {
        fontWeight: '300',
        color: '#808080',
    }
})